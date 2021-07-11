/* eslint-disable */
import { DenomTrace } from "../market/denomTrace";
import { OrderBook } from "../market/orderBook";
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "onomyprotocol.onomy.market";

/** GenesisState defines the market module's genesis state. */
export interface GenesisState {
  /** this line is used by starport scaffolding # genesis/proto/state */
  denomTraceList: DenomTrace[];
  /** this line is used by starport scaffolding # genesis/proto/stateField */
  orderBookList: OrderBook[];
  /** this line is used by starport scaffolding # genesis/proto/stateField */
  portId: string;
}

const baseGenesisState: object = { portId: "" };

export const GenesisState = {
  encode(message: GenesisState, writer: Writer = Writer.create()): Writer {
    for (const v of message.denomTraceList) {
      DenomTrace.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.orderBookList) {
      OrderBook.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.portId !== "") {
      writer.uint32(10).string(message.portId);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseGenesisState } as GenesisState;
    message.denomTraceList = [];
    message.orderBookList = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 3:
          message.denomTraceList.push(
            DenomTrace.decode(reader, reader.uint32())
          );
          break;
        case 2:
          message.orderBookList.push(OrderBook.decode(reader, reader.uint32()));
          break;
        case 1:
          message.portId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenesisState {
    const message = { ...baseGenesisState } as GenesisState;
    message.denomTraceList = [];
    message.orderBookList = [];
    if (object.denomTraceList !== undefined && object.denomTraceList !== null) {
      for (const e of object.denomTraceList) {
        message.denomTraceList.push(DenomTrace.fromJSON(e));
      }
    }
    if (object.orderBookList !== undefined && object.orderBookList !== null) {
      for (const e of object.orderBookList) {
        message.orderBookList.push(OrderBook.fromJSON(e));
      }
    }
    if (object.portId !== undefined && object.portId !== null) {
      message.portId = String(object.portId);
    } else {
      message.portId = "";
    }
    return message;
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    if (message.denomTraceList) {
      obj.denomTraceList = message.denomTraceList.map((e) =>
        e ? DenomTrace.toJSON(e) : undefined
      );
    } else {
      obj.denomTraceList = [];
    }
    if (message.orderBookList) {
      obj.orderBookList = message.orderBookList.map((e) =>
        e ? OrderBook.toJSON(e) : undefined
      );
    } else {
      obj.orderBookList = [];
    }
    message.portId !== undefined && (obj.portId = message.portId);
    return obj;
  },

  fromPartial(object: DeepPartial<GenesisState>): GenesisState {
    const message = { ...baseGenesisState } as GenesisState;
    message.denomTraceList = [];
    message.orderBookList = [];
    if (object.denomTraceList !== undefined && object.denomTraceList !== null) {
      for (const e of object.denomTraceList) {
        message.denomTraceList.push(DenomTrace.fromPartial(e));
      }
    }
    if (object.orderBookList !== undefined && object.orderBookList !== null) {
      for (const e of object.orderBookList) {
        message.orderBookList.push(OrderBook.fromPartial(e));
      }
    }
    if (object.portId !== undefined && object.portId !== null) {
      message.portId = object.portId;
    } else {
      message.portId = "";
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
