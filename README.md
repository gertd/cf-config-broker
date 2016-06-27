# cf-config-broker

Steps for installation:

1.  Deploy the broker
    cf push cf-config-broker -o gertd/cf-config-broker -m 32M

2.  Create the service broker
    cf create-service-broker cf-config-broker username password http://cf-config-broker.hcf.helion.io

3.  Enable public access to the services exposed
    cf enable-service-access mssql
    cf enable-service-access oradb

Using the cf-config-broker:

1.  Create service instance for each flavor
    cf create-service mssql default mssql-test
    cf create-service oradb default oradb-test

2.  Bind service to the php-env application
    cf bind-service php-env mssql-test
    cf bind-service php-env oradb-test

3.  Restage the application 
    cf restage php-env

Removing the cf-config-broker

1.  Unbind all application using services exposed by the service broker
    cf unbind-service php-env mssql-test
    cf unbind-service php-env oradb-test

2.  Delete the service instances
    cf delete-service mssql-test -f
    cf delete-service oradb-test -f

3.  Delete the service broker instance 
    cf delete-service-broker cf-config-broker -f

4.  Delete the service broker service (app)
    cf delete cf-config-broker -f
