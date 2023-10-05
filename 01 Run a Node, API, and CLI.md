<img width="1109" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/1385332f-cc9a-4301-89c8-d91c32829314">### Cosmos_sdk_study
코스모스 sdk 공부
```
$ mkdir cosmos
$ cd cosmos
$ git clone https://github.com/cosmos/cosmos-sdk
$ cd cosmos-sdk
```
<img width="1491" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/446b635f-b6da-44f7-849a-2dc2e848daf7">

```
$ git checkout v0.45.4
```

<img width="786" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/d50ff3cb-22f7-4abd-8e75-fc6ca7bea1cc">


```
$ make build
```

<img width="839" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/b88270af-1046-40cb-b9bd-251310a99710">

```
$ ./build/simd init demo --home ./private/.simapp --chain-id learning-chain-1
```

<img width="1775" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets  /33992354/507c7780-a886-443d-8f65-5977df4a6d85">

```
cat ./private/.simapp/config/genesis.json
```

<img width="1459" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/d8268aad-c7bc-410a-8a94-db13fd262a4d">

./build/simd keys list --home ./private/.simapp --keyring-backend test
키 리스트 확인

<img width="662" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/cb94ec10-734b-475f-a9e4-9a171ca526ed">

```
./build/simd keys add alice \
    --home ./private/.simapp \
    --keyring-backend test

grep -A 10 balances ./private/.simapp/config/genesis.json
```
Alice에게 키 생성 (이제 리스트에 Alice의 키가 나오게 된다)

<img width="1109" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/5f1c8d7a-f65c-482d-92eb-432f44218b87">

```
./build/simd add-genesis-account alice 100000000stake \
    --home ./private/.simapp \
    --keyring-backend test
```

<img width="1109" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/456a7acc-77b7-4bef-a088-3b46d5fca1d2">

```
./build/simd gentx alice 70000000stake \
    --home ./private/.simapp \
    --keyring-backend test \
    --chain-id learning-chain-1

genesis tx 생성
```

<img width="1109" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/d0fdc457-b3ad-4427-ad92-4fda84af3f40">

```
./build/simd start --home ./private/.simapp

블록체인 구동
```

<img width="1109" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/a4bcab0d-5d68-4e96-b39b-73dd7cf57643">

```
./build/simd query bank balances $alice

구동중인 블록체인에 query 날려 Alice의 balance 확인
```

<img width="1109" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/b5ebef82-b636-47fb-b789-cecd2ff16ebc">

```
./build/simd query bank balances $bob

Bob의 balance 확인
```

<img width="1109" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/f8964482-f288-4b6e-bc5c-59864eb38ac5">

```
./build/simd tx bank send $alice $bob 10stake \
    --home ./private/.simapp \
    --keyring-backend test \
    --chain-id learning-chain-1

Alice에서 Bob에게 10 토큰 전송
```

<img width="1109" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/f6d222a2-96c7-4613-9f58-35d63863c8bf">

```
./build/simd query tx 4D7C5F40FD7322E7B13581CAD637D96B902D6370986703ADC29F745DC77AC596

txHash를 쿼리하여 tx정보 받아오기
```

<img width="1278" alt="image" src="https://github.com/dik654/Cosmos_sdk_study/assets/33992354/85431319-b832-4ef4-8ce5-2deceb00ae4c">
