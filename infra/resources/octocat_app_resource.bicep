targetScope = 'subscription'
param epochTime string

resource containerAppGroup 'Microsoft.Resources/resourceGroups@2021-04-01' = {
  name: 'OctoCatRoasterApp'
  location: deployment().location
}

module containerAppResources 'container_app.bicep' = {
    name: 'OctocatRoasterContainerInnerApp'
    scope: containerAppGroup
    params: {
        epochTime: epochTime
    }
}

module frontendAppResource 'frontend_app.bicep' = {
    name: 'OctocatRoasterFrontendInnerApp'
    scope: containerAppGroup
}
