# RTI Shapes demo with TIG

Docker Compose file for running TIG (Telegraf, InfluxDB, and Garafana) with [RTI Shapes demo](https://www.rti.com/free-trial/shapes-demo).
This is to demonstrate [Telegraf plugin for RTI Connext DDS](https://www.rti.com/developers/rti-labs/telegraf-plugin-for-connext-dds). 
This requires installation of Docker and Docker Compose (above version 3).

## Usage

To persist data for InfluxDB and Grafana, you need to create directories:

    mkdir -p $HOME/docker/influxdb/data
    mkdir -p $HOME/docker/grafana/data; sudo chown 472:472 $HOME/docker/grafana/data

Run Docker Compose:
  
    docker-compose up -d
    
After running the Docker images, you can run the RTI Shapes demo and create Shapes publishers. 
Then, you can see Shapes data visualized in a Grafana dashboard (http://localhost:3000).

The baseline QoS setting used by the DataReaders in Telegraf is `Generic.KeepLastReliable`. 
The `Square` DataReader and Circle` DataReader` use the default QoS settings. 

The `Circle` DataReader uses a content-based filter (`x > 100`), so it will receive data only when the x coordinate is higher than 100. 

The `Triangle` DataReader uses `TRANSIENT_LOCAL_DURABILITY_QOS` and `KEEP_ALL_HISTORY_QOS`, so it will receive all historical data. 
