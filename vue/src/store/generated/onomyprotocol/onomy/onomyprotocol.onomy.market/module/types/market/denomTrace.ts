/* eslint-disable */
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "onomyprotocol.onomy.market";

export interface DenomTrace {
  index: string;
  port: string;
  channel: string;
  origin: string;
}

const baseDenomTrace: object = { index: "", port: "", channel: "", origin: "" };

export const DenomTrace = {
  encode(message: DenomTrace, writer: Writer = Writer.create()): Writer {
    if (message.index !== "") {
      writer.uint32(18).string(message.index);
    }
    if (message.port !== "") {
      writer.uint32(26).string(message.port);
    }
    if (message.channel !== "") {
      writer.uint32(34).string(message.channel);
    }
    if (message.origin !== "") {
      writer.uint32(42).string(message.origin);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): DenomTrace {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseDenomTrace } as DenomTrace;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.index = reader.string();
          break;
        case 3:
          message.port = reader.string();
          break;
        case 4:
          message.channel = reader.string();
          break;
        case 5:
          message.origin = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DenomTrace {
    const message = { ...baseDenomTrace } as DenomTrace;
    if (object.index !== undefined && object.index !== null) {
      message.index = String(object.index);
    } else {
      message.index = "";
    }
    if (object.port !== undefined && object.port !== null) {
      message.port = String(object.port);
    } else {
      message.port = "";
    }
    if (object.channel !== undefined && object.channel !== null) {
      message.channel = String(object.channel);
    } else {
      message.channel = "";
    }
    if (object.origin !== undefined && object.origin !== null) {
      message.origin = String(object.origin);
    } else {
      message.origin = "";
    }
    return message;
  },

  toJSON(message: DenomTrace): unknown {
    const obj: any = {};
    message.index !== undefined && (obj.index = message.index);
    message.port !== undefined && (obj.port = message.port);
    message.channel !== undefined && (obj.channel = message.channel);
    message.origin !== undefined && (obj.origin = message.origin);
    return obj;
  },

  fromPartial(object: DeepPartial<DenomTrace>): DenomTrace {
    const message = { ...baseDenomTrace } as DenomTrace;
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    } else {
      message.index = "";
    }
    if (object.port !== undefined && object.port !== null) {
      message.port = object.port;
    } else {
      message.port = "";
    }
    if (object.channel !== undefined && object.channel !== null) {
      message.channel = object.channel;
    } else {
      message.channel = "";
    }
    if (object.origin !== undefined && object.origin !== null) {
      message.origin = object.origin;
    } else {
      message.origin = "";
    }
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;
