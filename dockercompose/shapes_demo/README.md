# RTI Shapes demo with TIG

Docker-compose file for running TIG (Telegraf, InfluxDB, and Garafana) with [RTI Shapes demo](https://www.rti.com/free-trial/shapes-demo).
This is to demonstrate [Telegraf plugin for RTI Connext DDS](https://www.rti.com/developers/rti-labs/telegraf-plugin-for-connext-dds). 
This requires installation of Docker and Docker Compose (above version 3).

## Usage

Run Docker Compose:
  
    docker-compose up -d
    
After running the Docker images, you can run the RTI Shapes demo and create Shapes publishers. 
Then, you can see Shapes data visualized in a Grafana dashboard (http://localhost:3000).

The baseline QoS setting used by the DataReaders in Telegraf is `Generic.KeepLastReliable`. 
The `Square` DataReader uses the default QoS settings. 
The `Circle` DataReader uses the default QoS settings. It also defines a content-based filter (`x > 100`), so it will receive data only the x coordinate is higher than 100. 
The `Triangle` DataReader uses `TRANSIENT_LOCAL_DURABILITY_QOS` and `KEEP_ALL_HISTORY_QOS`. So it will receive all historical data from a DataWriter. 
