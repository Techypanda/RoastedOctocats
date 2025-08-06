resource reactStorageAccount 'Microsoft.Storage/storageAccounts@2024-01-01' = {
    name: 'octocatroaster'
    kind: 'StorageV2'
    location: resourceGroup().location
    properties: {
        allowBlobPublicAccess: true
    }
    sku: {
        name: 'Standard_LRS'
    }
}

resource reactBlobServices 'Microsoft.Storage/storageAccounts/blobServices@2024-01-01' = {
  name: 'default'
  parent: reactStorageAccount
}

resource reactStorageContainer 'Microsoft.Storage/storageAccounts/blobServices/containers@2024-01-01' = {
    parent: reactBlobServices
    name: 'www'
    properties: {
        publicAccess: 'Blob'
    }
}

resource reactCDN 'Microsoft.Cdn/profiles@2024-09-01' = {
    name: 'octocatroaster'
    location: resourceGroup().location
    sku: {
        name: 'Standard_Microsoft'
    }
}

resource reactEndpoint 'Microsoft.Cdn/profiles/endpoints@2024-09-01' = {
    parent: reactCDN
    location: resourceGroup().location
    name: 'octocatroaster'
    properties: {
        origins: [
            {
                name: 'octocatroaster'
                properties: {
                    enabled: true
                    hostName: replace(replace(reactStorageAccount.properties.primaryEndpoints.blob, 'https://', ''), 'net/', 'net')
                    httpPort: 80
                    httpsPort: 443
                    originHostHeader: replace(replace(reactStorageAccount.properties.primaryEndpoints.blob, 'https://', ''), 'net/', 'net')
                    priority: 1
                }
            }
        ]
        deliveryPolicy: {
            rules: [
                {
                    actions: [
                        {
                            name: 'CacheExpiration'
                            parameters: {
                                cacheBehavior: 'Override'
                                cacheDuration: '0.00:05:00'
                                cacheType: 'All'
                                typeName: 'DeliveryRuleCacheExpirationActionParameters'
                            }
                        }
                    ]
                    name: 'Always Cache expiration'
                    order: 0
                }
                {
                    conditions: [
                        {
                            name: 'UrlFileExtension'
                            parameters: {
                                operator: 'LessThan'
                                matchValues: [
                                    '1'
                                ]
                                typeName: 'DeliveryRuleUrlFileExtensionMatchConditionParameters'
                            }
                        }
                    ]
                    actions: [
                        {
                            name: 'UrlRewrite'
                            parameters: {
                                destination: '/index.html'
                                preserveUnmatchedPath: false
                                sourcePattern: '/'
                                typeName: 'DeliveryRuleUrlRewriteActionParameters'
                            }
                        }
                    ]
                    name: 'React Rewrite Rule'
                    order: 1
                }
            ]
        }
        originPath: '/www'
        optimizationType: 'GeneralWebDelivery'
    }
}
