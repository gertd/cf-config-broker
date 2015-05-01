# cf-config-broker

Steps for installation:

1.  git clone https://github.com/gertd/cf-config-broker/cf-config-broker.git
2.  cd cf-config-broker
3.  cp cf-config-broker.json new-broker-config.json
4.  edit the new-broker-config.json file
5.  deploy the broker using cf push, this will provide the network endpoint that can be used to bind the broker instance to
    helion push -n
6.  create service broker instance
    helion add-service-broker cf-config-broker
7.  update plan to become public
    helion update-service-plan --vendor mssql --public free

To use the broker:

1.  create service instance for each flavor
    helion create-service mssql mssql-db
2.  bind service to the application
    helion bind-service oracle-db go-env

