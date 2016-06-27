#!/bin/bash

cf unbind-service php-env mssql-test
cf unbind-service php-env oradb-test

cf delete-service mssql-test -f
cf delete-service oradb-test -f

cf delete-service-broker cf-config-broker -f

cf delete cf-config-broker -f

cf push cf-config-broker -o gertd/cf-config-broker -m 32M

cf create-service-broker cf-config-broker username password http://cf-config-broker.hcf.helion.io

cf enable-service-access mssql
cf enable-service-access oradb

cf create-service mssql default mssql-test
cf create-service oradb default oradb-test

cf bind-service php-env mssql-test
cf bind-service php-env oradb-test

cf restage php-env
