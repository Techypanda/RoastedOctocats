$rgName = "OctocatDeploymentRG"
$groupExists = az group exists -n $rgName
$epochTime = [int][double]::Parse((Get-Date -UFormat %s))
if ( $groupExists -eq "false" )
{
	Write-Warning "Resource Group For Bicep Deployment Doesnt Exist, Creating..."
	az group create --name $rgName --location "australiasoutheast"
}
echo "Building API Docker Image..."
Push-Location
Set-Location "../api_v2"
docker build . -t octocatroaster:$epochTime
docker build . -f Dockerfile.envoy -t octocatroasterenvoy:$epochTime
Pop-Location
echo "Pushing Docker Image To ACR"
docker tag octocatroaster:$epochTime octocatregistry.azurecr.io/octocatroaster:$epochTime
docker tag octocatroasterenvoy:$epochTime octocatregistry.azurecr.io/octocatroasterenvoy:$epochTime
az acr login --name octocatregistry.azurecr.io
docker push octocatregistry.azurecr.io/octocatroaster:$epochTime
docker push octocatregistry.azurecr.io/octocatroasterenvoy:$epochTime
echo "Building Frontend Code..."
Push-Location
Set-Location "../web"
npm run build
Pop-Location
echo "Now Deploying Main Infrastructure..."
az deployment group create --name MainInfrastructure --resource-group $rgName --parameters epochTime=$epochTime --template-file main.bicep --mode Complete
echo "Deploying frontend..."
az storage blob upload-batch --destination www --account-name octocatroaster --source ../web/dist --overwrite
Remove-Item "../web/dist" -Recurse -Force