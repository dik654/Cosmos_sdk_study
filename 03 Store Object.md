```
$ mkdir x/checkers/rules
$ curl https://raw.githubusercontent.com/batkinson/checkers-go/a09daeb1548dd4cc0145d87c8da3ed2ea33a62e3/checkers/checkers.go | sed 's/package checkers/package rules/' > x/checkers/rules/checkers.go

서버 파일 없이 게임 파일만 받아서 rules에 저장
```

<img width="1780" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/62a0f1c2-80c8-49ab-8f12-119528123f26">

```
### 저장해야할 Object
- 흑돌 address string
- 붉은돌 address string
- 게임판 상태 string
- 게임 턴
```

```
각각의 한판한판 게임들을 미래에도 특정지을 수 있어야하기에 SystemInfo안에 counter를 생성
$ ignite scaffold single systemInfo nextId:uint \
    --module checkers \
    --no-message

--no-message를 두지않으면 SystemInfo를 덮어쓰는 sdk.Msg를 생성하므로 붙여줘야한다
```

<img width="1775" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/b4a60bc2-8be0-468c-ac33-9efc2e5526dd">

```
$ ignite scaffold map storedGame board turn black red \
    --index index \
    --module checkers \
    --no-message

이제 위에서 생성한 count를 index로 Object를 저장하기위한 StoredGame을 생성한다
SystemInfo처럼 Msg로 인해 object들이 생성되거나 수정되는 것을 방지하기 위해 --no-message를 붙인다

map은 항상 key로 string 타입을 가진다
```

<img width="1775" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/468845f8-4000-4ada-aeac-27a6c0b697ba">

```
const (
    SystemInfoKey = "SystemInfo-value-"
)

const (
    StoredGameKeyPrefix = "StoredGame/value/"
)

object에 접근하기위한 prefix
```

<img width="1775" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/f217914c-fde0-4e98-b765-ffab02fe38f4">

```
store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StoredGameKeyPrefix))

변수를 저장하는 store 접근이 용이하도록 prefix설정되는 부분

b := store.Get(types.StoredGameKey(
    index,
))

인덱스로 game store에 접근할 수 있도록 한다
```

<img width="1764" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/8d76b1b1-0d5e-4d7a-b8f1-ce840cc0b422">

```
컴파일 전 protobuf object는 아래와 같이 message 형식을 갖는다

message SystemInfo {
    uint64 nextId = 1;
}

message StoredGame {
    string index = 1;
    string board = 2;
    string turn = 3;
    string black = 4;
    string red = 5;
}

컴파일 후는 아래와 같다

type SystemInfo struct {
    NextId uint64 `protobuf:"varint,1,opt,name=nextId,proto3" json:"nextId,omitempty"`
}

type StoredGame struct {
    Index string `protobuf:"bytes,1,opt,name=index,proto3" json:"index,omitempty"`
    Board string `protobuf:"bytes,2,opt,name=board,proto3" json:"board,omitempty"`
    Turn  string `protobuf:"bytes,3,opt,name=turn,proto3" json:"turn,omitempty"`
    Black string `protobuf:"bytes,4,opt,name=black,proto3" json:"black,omitempty"`
    Red   string `protobuf:"bytes,5,opt,name=red,proto3" json:"red,omitempty"`
}

```
컴파일 전 protobuf의 genesis state도

import "checkers/system_info.proto";
import "checkers/stored_game.proto";

message GenesisState {
    ...
    SystemInfo systemInfo = 2;
    repeated StoredGame storedGameList = 3 [(gogoproto.nullable) = false];
}

아래와 같이 컴파일된다

type GenesisState struct {
    Params         Params       `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
    SystemInfo     *SystemInfo  `protobuf:"bytes,2,opt,name=systemInfo,proto3" json:"systemInfo,omitempty"`
    StoredGameList []StoredGame `protobuf:"bytes,3,rep,name=storedGameList,proto3" json:"storedGameList"`
}
```

```
StoredGame의 쿼리 object들

message QueryGetSystemInfoRequest {}

message QueryGetSystemInfoResponse {
    SystemInfo SystemInfo = 1 [(gogoproto.nullable) = false];
}

message QueryGetStoredGameRequest {
    string index = 1;
}

message QueryGetStoredGameResponse {
    StoredGame StoredGame = 1 [(gogoproto.nullable) = false];
}

message QueryAllStoredGameRequest {
    cosmos.base.query.v1beta1.PageRequest pagination = 1;
}

message QueryAllStoredGameResponse {
    repeated StoredGame StoredGame = 1 [(gogoproto.nullable) = false];
    cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

```
null을 허용하게 기본 설정되어있는 nextId에 초기값과 nullable을 끄기 위한 설정
```

<img width="1764" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/7add95e6-cc44-480d-b1ab-92d0d2767d3e">

<img width="1764" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/e8f93459-a4ef-467c-8df4-bba8403247f9">

<img width="1764" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/ef828c3f-6cec-40c3-a23d-ca16fa430344">

<img width="1764" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/5a3c062b-17d9-44a7-8a37-ffbd8598f6ad">

<img width="1764" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/ad2e6acc-b183-4507-9f21-c707ae99f574">

이후 과정은 import하는 repo가 삭제되어 빌드가 되지않아 실제로 실행해보진 못했다
<img width="1764" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/7e2b47f6-be25-4ab6-9848-fa2319b32f7e">


