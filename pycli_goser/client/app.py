import grpc

import mamushi_pb2
import mamushi_pb2_grpc


def run():
    msg = input("何か入力してください")
    with grpc.insecure_channel("localhost:50051") as ch:
        stub = mamushi_pb2_grpc.MamushiStub(ch)
        stub.PushMsg(mamushi_pb2.MsgStruct(message=msg))


if __name__ == '__main__':
    run()
