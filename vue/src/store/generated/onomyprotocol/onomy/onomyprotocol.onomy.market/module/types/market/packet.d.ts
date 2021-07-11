import { Writer, Reader } from "protobufjs/minimal";
export declare const protobufPackage = "onomyprotocol.onomy.market";
export interface MarketPacketData {
    noData: NoData | undefined;
    /** this line is used by starport scaffolding # ibc/packet/proto/field */
    createOrderPacket: CreateOrderPacketData | undefined;
    /** this line is used by starport scaffolding # ibc/packet/proto/field/number */
    createPairPacket: CreatePairPacketData | undefined;
}
export interface NoData {
}
/**
 * this line is used by starport scaffolding # ibc/packet/proto/message
 * CreateOrderPacketData defines a struct for the packet payload
 */
export interface CreateOrderPacketData {
    amountDenom: string;
    amount: number;
    sourceCoin: string;
    targetCoin: string;
    exchRateDenom: string;
    exchRate: string;
}
/** CreateOrderPacketAck defines a struct for the packet acknowledgment */
export interface CreateOrderPacketAck {
}
/** CreatePairPacketData defines a struct for the packet payload */
export interface CreatePairPacketData {
    sourceDenom: string;
    targetDenom: string;
}
/** CreatePairPacketAck defines a struct for the packet acknowledgment */
export interface CreatePairPacketAck {
}
export declare const MarketPacketData: {
    encode(message: MarketPacketData, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): MarketPacketData;
    fromJSON(object: any): MarketPacketData;
    toJSON(message: MarketPacketData): unknown;
    fromPartial(object: DeepPartial<MarketPacketData>): MarketPacketData;
};
export declare const NoData: {
    encode(_: NoData, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): NoData;
    fromJSON(_: any): NoData;
    toJSON(_: NoData): unknown;
    fromPartial(_: DeepPartial<NoData>): NoData;
};
export declare const CreateOrderPacketData: {
    encode(message: CreateOrderPacketData, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): CreateOrderPacketData;
    fromJSON(object: any): CreateOrderPacketData;
    toJSON(message: CreateOrderPacketData): unknown;
    fromPartial(object: DeepPartial<CreateOrderPacketData>): CreateOrderPacketData;
};
export declare const CreateOrderPacketAck: {
    encode(_: CreateOrderPacketAck, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): CreateOrderPacketAck;
    fromJSON(_: any): CreateOrderPacketAck;
    toJSON(_: CreateOrderPacketAck): unknown;
    fromPartial(_: DeepPartial<CreateOrderPacketAck>): CreateOrderPacketAck;
};
export declare const CreatePairPacketData: {
    encode(message: CreatePairPacketData, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): CreatePairPacketData;
    fromJSON(object: any): CreatePairPacketData;
    toJSON(message: CreatePairPacketData): unknown;
    fromPartial(object: DeepPartial<CreatePairPacketData>): CreatePairPacketData;
};
export declare const CreatePairPacketAck: {
    encode(_: CreatePairPacketAck, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): CreatePairPacketAck;
    fromJSON(_: any): CreatePairPacketAck;
    toJSON(_: CreatePairPacketAck): unknown;
    fromPartial(_: DeepPartial<CreatePairPacketAck>): CreatePairPacketAck;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
