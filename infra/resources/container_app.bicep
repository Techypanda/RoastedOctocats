param epochTime string

resource identity 'Microsoft.ManagedIdentity/userAssignedIdentities@2025-01-31-preview' = {
  name: 'OctocatContainerIdentity'
  location: resourceGroup().location
}

// TODO: Deploy the OpenAI model in this RG and update permissions

resource octocatregistry 'Microsoft.ContainerRegistry/registries@2024-11-01-preview' = {
    name: 'octocatregistry'
    location: resourceGroup().location
    sku: {
        name: 'Basic'
    }
    properties: {
        adminUserEnabled: false
        anonymousPullEnabled: false
        networkRuleBypassOptions: 'AzureServices'
    }
}

resource roleAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(resourceGroup().id, 'octocatregistry', 'AcrPullTestUserAssigned')
  properties: {
    principalId: identity.properties.principalId  
    principalType: 'ServicePrincipal'
    // acrPullDefinitionId has a value of 7f951dda-4ed3-4680-a7ca-43fe172d538d
    roleDefinitionId: resourceId('Microsoft.Authorization/roleDefinitions', '7f951dda-4ed3-4680-a7ca-43fe172d538d')
  }
}

resource aiRoleAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(resourceGroup().id, 'AIContributor')
  properties: {
    principalId: identity.properties.principalId  
    principalType: 'ServicePrincipal'
    // aicontrib has a value of a001fd3d-188f-4b5d-821b-7da978bf7442
    roleDefinitionId: resourceId('Microsoft.Authorization/roleDefinitions', 'a001fd3d-188f-4b5d-821b-7da978bf7442')
  }
}

resource law 'Microsoft.OperationalInsights/workspaces@2023-09-01' = {
  name: 'octocatLogAnalyticsWorkspace'
  location: resourceGroup().location
  properties: any({
    retentionInDays: 30
    features: {
      searchVersion: 1
    }
    sku: {
      name: 'PerGB2018'
    }
  })
}

resource env 'Microsoft.App/managedEnvironments@2024-10-02-preview' = {
  name: 'octocatContainerAppEnvironment'
  location: resourceGroup().location
  properties: {
    appLogsConfiguration: {
      destination: 'log-analytics'
      logAnalyticsConfiguration: {
        customerId: law.properties.customerId
        sharedKey: law.listKeys().primarySharedKey
      }
    }
  }
}

resource dummyContainerApp 'Microsoft.App/containerApps@2024-10-02-preview' = {
    name: 'octocatroastercontainerapp'
    kind: 'containerapp'
    location: resourceGroup().location
    properties: {
        managedEnvironmentId: env.id
        configuration: {
            activeRevisionsMode: 'Single'
            registries: [
                {
                    identity: identity.id
                    server: octocatregistry.properties.loginServer
                }
            ]
            ingress: {
                targetPort: 8081
                external: true
                corsPolicy: {
                    allowCredentials: true
                    allowedOrigins: ['http://localhost:8081'] // , 'https://jonwright-azure-react.azureedge.net']
                }
            }
        }
        template: {
            containers: [
                {
                    name: 'octoroastergrpc'
                    image: '${octocatregistry.properties.loginServer}/octocatroaster:${epochTime}'
                    env: [
                      {
                        name: 'ADDRESS'
                        value: '0.0.0.0'
                      }
                      {
                        name: 'PORT'
                        value: '8080'
                      }
                      {
                        name: 'ENABLE_DEBUG'
                        value: 'true'
                      }
                      {
                        name: 'APPLICATION_ENVIRONMENT'
                        value: 'local'
                      }
                      {
                        name: 'BASEPATH'
                        value: '/app'
                      }
                      {
                        name: 'AZURE_CLIENT_ID'
                        value: identity.properties.clientId
                      }
                    ]
                    resources: {
                        cpu: '0.25'
                        memory: '0.5Gi'
                    }
                }
                {
                    name: 'octoroasterenvoy'
                    image: '${octocatregistry.properties.loginServer}/octocatroasterenvoy:${epochTime}'
                    resources: {
                        cpu: '0.25'
                        memory: '0.5Gi'
                    }
                    command: ['/usr/local/bin/envoy', '-c', '/etc/envoy/envoy_acs.yaml', '-l', 'trace', '--log-path', '/tmp/envoy_info.log']
                }
            ]
            scale: {
                maxReplicas: 1
                // minReplicas: 1
            }
        }
    }
    identity: {
        type: 'UserAssigned'
        userAssignedIdentities: {
            '${identity.id}': {}
        }
    }
}
