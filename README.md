A plugin to Drone Flyway runner DB migration support plugin.

# Usage

The following settings changes this plugin's behavior.

| **Plugin parameter** | **Description**                                                                                                   |
|------------------------|-------------------------------------------------------------------------------------------------------------------|
| **driver_path**        | Path to JDBC driver jar files when user wants a specific set of drivers to be used. e.g "/path/to/driver1.jar:/path/to/driver2.jar" |
| **flyway_command**     | The Flyway operation supports actions such as migrate, clean, baseline, or repair                                |
| **locations**          | Specifies the path to the Flyway migration files. This is where Flyway looks for SQL scripts. e.g /opt/web-app-3/migrations |
| **command_line_args**  | Additional arguments passed to Flyway. e.g “-Dflyway.schemas=public“ or “-X” etc …                                |
| **url**                | JDBC connection URL for the target database. e.g for mysql jdbc:mysql://4.20.19.21:3306/flyway_test               |
| **username**           | User name of the DB user                                                                                         |
| **password**           | Password of the DB user                                                                                          |

<br>

Below is an example `.drone.yml` that uses this plugin.

```yaml
- step:
    identifier: flywayrunner7f66cc
    name: flywayrunner
    spec:
      image: plugins/drone-flyway-runner
      settings:
        command_line_args: -X
        flyway_command: migrate
        locations: /opt/hns/harness-plugins/flyway-test-files/migration_files
        password: <+input>
        url: jdbc:mysql://43.204.190.241:3306/flyway_test
        username: <+input>
    timeout: ""
    type: Plugin
```

# Building

Build the plugin binary:

```text
scripts/build.sh
```

Build the plugin image:

```text
docker build -t plugins/drone-flyway-runner -f docker/Dockerfile .
```

# Testing

Execute the plugin from your current working directory:

### Migrate command example

```bash
docker run --rm -e \  
  -e PLUGIN_FLYWAY_COMMAND=migrate \
  -e PLUGIN_USERNAME=hns_usr \
  -e PLUGIN_PASSWORD=sksi$k89 \
  -e PLUGIN_URL=jdbc:mysql://23.12.7.8:3306/flyway_db \
  -e PLUGIN_LOCATIONS=/harness/db-migrations  
  plugins/drone-flyway-runner
```
### Migrate command with config file example

```bash
docker run --rm -e \  
  -e PLUGIN_FLYWAY_COMMAND=migrate \
  -e PLUGIN_COMMAND_LINE_ARGS='-configFiles=/harness/hns/test-resources/flyway/config1/flyway.conf'

```

### Migrate command with JDBC driver example

```bash
docker run --rm -e \  
  -e PLUGIN_FLYWAY_COMMAND=migrate \
  -e PLUGIN_USERNAME=hns_usr \
  -e PLUGIN_PASSWORD=sksi$k89 \
  -e PLUGIN_URL=jdbc:mysql://23.12.7.8:3306/flyway_db \
  -e PLUGIN_LOCATIONS=/harness/db-migrations \
  -e PLUGIN_DRIVER_PATH=/harness/drivers/mysql-connector-java-8.0.23.jar \  
  plugins/drone-flyway-runner
```

### Repair command with command line args example

```bash
docker run --rm -e \  
  -e PLUGIN_FLYWAY_COMMAND=repair \
  -e PLUGIN_USERNAME=hnstest03 \
  -e PLUGIN_PASSWORD=3cbc98835323 \
  -e PLUGIN_URL=jdbc:mysql://23.12.7.8:3306/flyway_db \
  -e PLUGIN_LOCATIONS=/harness/db-migrations \
  -e PLUGIN_COMMAND_LINE_ARGS='-connectRetries=3 -cleanDisabled=false -validateOnMigrate=true -baselineOnMigrate=true'
```
