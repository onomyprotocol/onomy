import { Reader, Writer } from "protobufjs/minimal";
import { DenomTrace } from "../market/denomTrace";
import { PageRequest, PageResponse } from "../cosmos/base/query/v1beta1/pagination";
import { OrderBook } from "../market/orderBook";
export declare const protobufPackage = "onomyprotocol.onomy.market";
/** this line is used by starport scaffolding # 3 */
export interface QueryGetDenomTraceRequest {
    index: string;
}
export interface QueryGetDenomTraceResponse {
    DenomTrace: DenomTrace | undefined;
}
export interface QueryAllDenomTraceRequest {
    pagination: PageRequest | undefined;
}
export interface QueryAllDenomTraceResponse {
    DenomTrace: DenomTrace[];
    pagination: PageResponse | undefined;
}
export interface QueryGetOrderBookRequest {
    index: string;
}
export interface QueryGetOrderBookResponse {
    OrderBook: OrderBook | undefined;
}
export interface QueryAllOrderBookRequest {
    pagination: PageRequest | undefined;
}
export interface QueryAllOrderBookResponse {
    OrderBook: OrderBook[];
    pagination: PageResponse | undefined;
}
export declare const QueryGetDenomTraceRequest: {
    encode(message: QueryGetDenomTraceRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryGetDenomTraceRequest;
    fromJSON(object: any): QueryGetDenomTraceRequest;
    toJSON(message: QueryGetDenomTraceRequest): unknown;
    fromPartial(object: DeepPartial<QueryGetDenomTraceRequest>): QueryGetDenomTraceRequest;
};
export declare const QueryGetDenomTraceResponse: {
    encode(message: QueryGetDenomTraceResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryGetDenomTraceResponse;
    fromJSON(object: any): QueryGetDenomTraceResponse;
    toJSON(message: QueryGetDenomTraceResponse): unknown;
    fromPartial(object: DeepPartial<QueryGetDenomTraceResponse>): QueryGetDenomTraceResponse;
};
export declare const QueryAllDenomTraceRequest: {
    encode(message: QueryAllDenomTraceRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryAllDenomTraceRequest;
    fromJSON(object: any): QueryAllDenomTraceRequest;
    toJSON(message: QueryAllDenomTraceRequest): unknown;
    fromPartial(object: DeepPartial<QueryAllDenomTraceRequest>): QueryAllDenomTraceRequest;
};
export declare const QueryAllDenomTraceResponse: {
    encode(message: QueryAllDenomTraceResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryAllDenomTraceResponse;
    fromJSON(object: any): QueryAllDenomTraceResponse;
    toJSON(message: QueryAllDenomTraceResponse): unknown;
    fromPartial(object: DeepPartial<QueryAllDenomTraceResponse>): QueryAllDenomTraceResponse;
};
export declare const QueryGetOrderBookRequest: {
    encode(message: QueryGetOrderBookRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryGetOrderBookRequest;
    fromJSON(object: any): QueryGetOrderBookRequest;
    toJSON(message: QueryGetOrderBookRequest): unknown;
    fromPartial(object: DeepPartial<QueryGetOrderBookRequest>): QueryGetOrderBookRequest;
};
export declare const QueryGetOrderBookResponse: {
    encode(message: QueryGetOrderBookResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryGetOrderBookResponse;
    fromJSON(object: any): QueryGetOrderBookResponse;
    toJSON(message: QueryGetOrderBookResponse): unknown;
    fromPartial(object: DeepPartial<QueryGetOrderBookResponse>): QueryGetOrderBookResponse;
};
export declare const QueryAllOrderBookRequest: {
    encode(message: QueryAllOrderBookRequest, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryAllOrderBookRequest;
    fromJSON(object: any): QueryAllOrderBookRequest;
    toJSON(message: QueryAllOrderBookRequest): unknown;
    fromPartial(object: DeepPartial<QueryAllOrderBookRequest>): QueryAllOrderBookRequest;
};
export declare const QueryAllOrderBookResponse: {
    encode(message: QueryAllOrderBookResponse, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): QueryAllOrderBookResponse;
    fromJSON(object: any): QueryAllOrderBookResponse;
    toJSON(message: QueryAllOrderBookResponse): unknown;
    fromPartial(object: DeepPartial<QueryAllOrderBookResponse>): QueryAllOrderBookResponse;
};
/** Query defines the gRPC querier service. */
export interface Query {
    /** this line is used by starport scaffolding # 2 */
    DenomTrace(request: QueryGetDenomTraceRequest): Promise<QueryGetDenomTraceResponse>;
    DenomTraceAll(request: QueryAllDenomTraceRequest): Promise<QueryAllDenomTraceResponse>;
    OrderBook(request: QueryGetOrderBookRequest): Promise<QueryGetOrderBookResponse>;
    OrderBookAll(request: QueryAllOrderBookRequest): Promise<QueryAllOrderBookResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    DenomTrace(request: QueryGetDenomTraceRequest): Promise<QueryGetDenomTraceResponse>;
    DenomTraceAll(request: QueryAllDenomTraceRequest): Promise<QueryAllDenomTraceResponse>;
    OrderBook(request: QueryGetOrderBookRequest): Promise<QueryGetOrderBookResponse>;
    OrderBookAll(request: QueryAllOrderBookRequest): Promise<QueryAllOrderBookResponse>;
}
interface Rpc {
    request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
