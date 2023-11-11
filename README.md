# bellis(in progress)

The distributed cloud monitoring and alarm system adopts a design with clear division of labor between modules. The modules communicate through GRPC and message queue (Redis), and the sampled data is stored through InfluxDB.

| Home page                                                                                              | Add monitoring target                                                                                  |
| ------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| ![image](https://github.com/bellis-daemon/bellis/assets/55825043/d79e7d3c-f01f-46e6-9bbf-22999de58ae6) | ![image](https://github.com/bellis-daemon/bellis/assets/55825043/89182639-126c-4762-abb0-cf1d3b22ae11) |

## Configure

### ETCD

| key                     | description                                 |
| ----------------------- | ------------------------------------------- |
| `influxdb_token`        | api token for influxdb.                     |
| `telegram_bot_token`    | telegram bot token generated via @BotFather |
| `tencent_smtp_password` | tencent cloud ses service smtp password     |
