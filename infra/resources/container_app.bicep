param epochTime string

resource identity 'Microsoft.ManagedIdentity/userAssignedIdentities@2025-01-31-preview' = {
  name: 'OctocatContainerIdentity'
  location: resourceGroup().location
}

resource tableAccount 'Microsoft.DocumentDB/databaseAccounts@2025-05-01-preview' = {
  name: 'octocatroasterdb'
  location: resourceGroup().location
  kind: 'GlobalDocumentDB'
  properties: {
    capabilities: [
      {
        name: 'EnableTable'
      }
    ]
    capacityMode: 'Serverless'
    databaseAccountOfferType: 'Standard'
    locations: [
      {
        failoverPriority: 0
        isZoneRedundant: false
        locationName: 'australiaeast'
      }
    ]
  }
}

resource table 'Microsoft.DocumentDB/databaseAccounts/tables@2025-05-01-preview' = {
  parent: tableAccount
  name: 'octoroastertable'
  location: resourceGroup().location
  properties: {
    resource: {
      id: 'octoroastertable'
      createMode: 'Default'
    }
  }
}

resource dbAssignment 'Microsoft.DocumentDB/databaseAccounts/sqlRoleAssignments@2024-05-15' = {
  name: guid(resourceGroup().id, 'octocatcontainer', 'DBAccess')
  parent: tableAccount
  properties: {
    principalId: identity.properties.principalId  
    // https://learn.microsoft.com/en-us/azure/cosmos-db/nosql/how-to-grant-data-plane-access?tabs=built-in-definition%2Ccsharp&pivots=azure-interface-bicep#permission-model
    roleDefinitionId: '/subscriptions/66a28eda-a79b-4181-b86b-3bac8e334da1/resourceGroups/OctoCatRoasterApp/providers/Microsoft.DocumentDB/databaseAccounts/octocatroasterdb/sqlRoleDefinitions/00000000-0000-0000-0000-000000000002'
    scope: tableAccount.id
  }
}
resource adminDBAssignment 'Microsoft.DocumentDB/databaseAccounts/sqlRoleAssignments@2024-05-15' = {
  name: guid(resourceGroup().id, 'octocatcontainer', 'AdminDBAccess')
  parent: tableAccount
  properties: {
    // Jonathan object id
    principalId: 'eb4d8a57-1f82-4ee2-8aa1-11e6894602e4'
    // https://learn.microsoft.com/en-us/azure/cosmos-db/nosql/how-to-grant-data-plane-access?tabs=built-in-definition%2Ccsharp&pivots=azure-interface-bicep#permission-model
    roleDefinitionId: '/subscriptions/66a28eda-a79b-4181-b86b-3bac8e334da1/resourceGroups/OctoCatRoasterApp/providers/Microsoft.DocumentDB/databaseAccounts/octocatroasterdb/sqlRoleDefinitions/00000000-0000-0000-0000-000000000002'
    scope: tableAccount.id
  }
}
resource localManagedIdentityDBAssignment 'Microsoft.DocumentDB/databaseAccounts/sqlRoleAssignments@2024-05-15' = {
  name: guid(resourceGroup().id, 'octocatcontainer', 'LocalDBAccess')
  parent: tableAccount
  properties: {
    // techypanda-local-dev
    principalId: 'e8bf2f09-8358-446c-9a5e-5ed365f09ff4'
    // https://learn.microsoft.com/en-us/azure/cosmos-db/nosql/how-to-grant-data-plane-access?tabs=built-in-definition%2Ccsharp&pivots=azure-interface-bicep#permission-model
    roleDefinitionId: '/subscriptions/66a28eda-a79b-4181-b86b-3bac8e334da1/resourceGroups/OctoCatRoasterApp/providers/Microsoft.DocumentDB/databaseAccounts/octocatroasterdb/sqlRoleDefinitions/00000000-0000-0000-0000-000000000002'
    scope: tableAccount.id
  }
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
