import { Order } from "../market/order";
import { Writer, Reader } from "protobufjs/minimal";
export declare const protobufPackage = "onomyprotocol.onomy.market";
export interface OrderBook {
    index: string;
    orderIDTrack: number;
    amountDenom: string;
    exchRateDenom: string;
    orders: Order[];
}
export declare const OrderBook: {
    encode(message: OrderBook, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): OrderBook;
    fromJSON(object: any): OrderBook;
    toJSON(message: OrderBook): unknown;
    fromPartial(object: DeepPartial<OrderBook>): OrderBook;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
