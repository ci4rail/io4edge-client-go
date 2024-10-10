```mermaid
classDiagram
    motionsensor_Client *-- functionblock_Client
    functionblock_Client *-- pbchannelclient_Client
    core_If ..|> pbcore_Client
    core_If ..|> restcore_Client
    pbcore_Client *-- pbchannelclient_Client
    pbchannelclient_Client *-- pbchannelclient_Channel
    pbchannelclient_ChannelIf ..|> pbchannelclient_Channel
    pbchannelclient_Channel  *-- transport_FramedStream
    transport_MsgStream  ..|> transport_FramedStream

class motionsensor_Client {
    - functionblock_Client fbClient
}

class functionblock_Client {
    - client_Client funcClient
    + UploadConfiguration
    + FunctionControlSet
    ...
}

class pbchannelclient_Client {
    + Channel ch
    + FunctionInfo FuncInfo
}

class core_If {
    + IdentifyFirmware
    + LoadFirmware
    + IdentifyHardware
    ...
}

class pbcore_Client {
    - client.Client funcClient
}

class pbchannelclient_Channel {
    - transport.MsgStream ms
}

class pbchannelclient_ChannelIf {
    Operates on protobuf messages
    + WriteMessage
    + ReadMessage
}

class transport_FramedStream {
    + transport.Transport Trans
}

class transport_MsgStream {
    Operates on Byte Streams
    + ReadMsg
    + WriteMsg
    + Close
}

class zeroconfservice {
    + ServiceObserver
    + GetServiceInfo
}

```

New structure

```
cmd/
pkg/
    protobufcom/
        common/
            channel
            functionblock
        core
        functionblockclients/
            analogintypea
            binaryiotypea
        tracelet
    restcom/
        core/
    zeroconfservice/
    core/
    transport
    server
```
