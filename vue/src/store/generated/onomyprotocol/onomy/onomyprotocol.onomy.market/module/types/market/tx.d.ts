import { Reader, Writer } from "protobufjs/minimal";
export declare const protobufPackage = "onomyprotocol.onomy.market";
/** this line is used by starport scaffolding # proto/tx/message */
export interface MsgCancelOrder {
    creator: string;
    port: string;
    channel: string;
    amountDenom: string;
    exchRateDenom: string;
    orderID: number;
}
export interface MsgCancelOrderResponse {
}
export interface MsgSendCreateOrder {
    sender: string;
    port: string;
    channelID: string;
    timeoutTimestamp: number;
    amountDenom: string;
    amount: number;
    sourceCoin: string;
    targetCoin: string;
    exchRateDenom: string;
    exchRate: string;
}
export interface MsgSendCreateOrderResponse {
}
export interface MsgSendCreatePair {
    sender: string;
    port: string;
    channelID: string;
    timeoutTimestamp: number;
    sourceDenom: string;
    targetDenom: string;
}
export interface MsgSendCreatePairResponse {
}
export declare const MsgCancelOrder: {
    encode(message: MsgCancelOrder, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgCancelOrder;
    fromJSON(object: any): MsgCancelOrder;
    toJSON(message: MsgCancelOrder): unknown;
    fromPartial(object: DeepPartial<MsgCancelOrder>): MsgCancelOrder;
};
export declare const MsgCancelOrderResponse: {
    encode(_: MsgCancelOrderResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgCancelOrderResponse;
    fromJSON(_: any): MsgCancelOrderResponse;
    toJSON(_: MsgCancelOrderResponse): unknown;
    fromPartial(_: DeepPartial<MsgCancelOrderResponse>): MsgCancelOrderResponse;
};
export declare const MsgSendCreateOrder: {
    encode(message: MsgSendCreateOrder, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgSendCreateOrder;
    fromJSON(object: any): MsgSendCreateOrder;
    toJSON(message: MsgSendCreateOrder): unknown;
    fromPartial(object: DeepPartial<MsgSendCreateOrder>): MsgSendCreateOrder;
};
export declare const MsgSendCreateOrderResponse: {
    encode(_: MsgSendCreateOrderResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgSendCreateOrderResponse;
    fromJSON(_: any): MsgSendCreateOrderResponse;
    toJSON(_: MsgSendCreateOrderResponse): unknown;
    fromPartial(_: DeepPartial<MsgSendCreateOrderResponse>): MsgSendCreateOrderResponse;
};
export declare const MsgSendCreatePair: {
    encode(message: MsgSendCreatePair, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgSendCreatePair;
    fromJSON(object: any): MsgSendCreatePair;
    toJSON(message: MsgSendCreatePair): unknown;
    fromPartial(object: DeepPartial<MsgSendCreatePair>): MsgSendCreatePair;
};
export declare const MsgSendCreatePairResponse: {
    encode(_: MsgSendCreatePairResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MsgSendCreatePairResponse;
    fromJSON(_: any): MsgSendCreatePairResponse;
    toJSON(_: MsgSendCreatePairResponse): unknown;
    fromPartial(_: DeepPartial<MsgSendCreatePairResponse>): MsgSendCreatePairResponse;
};
/** Msg defines the Msg service. */
export interface Msg {
    /** this line is used by starport scaffolding # proto/tx/rpc */
    CancelOrder(request: MsgCancelOrder): Promise<MsgCancelOrderResponse>;
    SendCreateOrder(request: MsgSendCreateOrder): Promise<MsgSendCreateOrderResponse>;
    SendCreatePair(request: MsgSendCreatePair): Promise<MsgSendCreatePairResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    CancelOrder(request: MsgCancelOrder): Promise<MsgCancelOrderResponse>;
    SendCreateOrder(request: MsgSendCreateOrder): Promise<MsgSendCreateOrderResponse>;
    SendCreatePair(request: MsgSendCreatePair): Promise<MsgSendCreatePairResponse>;
}
interface Rpc {
    request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
