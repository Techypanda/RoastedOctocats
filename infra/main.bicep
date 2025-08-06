param epochTime string

module octocatContainerApp 'resources/octocat_app_resource.bicep' = {
    name: 'octocatContainerApp'
    scope: subscription()
    params: {
        epochTime: epochTime
    }
}
