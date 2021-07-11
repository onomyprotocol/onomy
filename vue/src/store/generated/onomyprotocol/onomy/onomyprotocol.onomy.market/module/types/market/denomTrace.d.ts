import { Writer, Reader } from "protobufjs/minimal";
export declare const protobufPackage = "onomyprotocol.onomy.market";
export interface DenomTrace {
    index: string;
    port: string;
    channel: string;
    origin: string;
}
export declare const DenomTrace: {
    encode(message: DenomTrace, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): DenomTrace;
    fromJSON(object: any): DenomTrace;
    toJSON(message: DenomTrace): unknown;
    fromPartial(object: DeepPartial<DenomTrace>): DenomTrace;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
