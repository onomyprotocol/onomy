/* eslint-disable */
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "onomyprotocol.onomy.market";

export interface MarketPacketData {
  noData: NoData | undefined;
  /** this line is used by starport scaffolding # ibc/packet/proto/field */
  createOrderPacket: CreateOrderPacketData | undefined;
  /** this line is used by starport scaffolding # ibc/packet/proto/field/number */
  createPairPacket: CreatePairPacketData | undefined;
}

export interface NoData {}

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
export interface CreateOrderPacketAck {}

/** CreatePairPacketData defines a struct for the packet payload */
export interface CreatePairPacketData {
  sourceDenom: string;
  targetDenom: string;
}

/** CreatePairPacketAck defines a struct for the packet acknowledgment */
export interface CreatePairPacketAck {}

const baseMarketPacketData: object = {};

export const MarketPacketData = {
  encode(message: MarketPacketData, writer: Writer = Writer.create()): Writer {
    if (message.noData !== undefined) {
      NoData.encode(message.noData, writer.uint32(10).fork()).ldelim();
    }
    if (message.createOrderPacket !== undefined) {
      CreateOrderPacketData.encode(
        message.createOrderPacket,
        writer.uint32(26).fork()
      ).ldelim();
    }
    if (message.createPairPacket !== undefined) {
      CreatePairPacketData.encode(
        message.createPairPacket,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MarketPacketData {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMarketPacketData } as MarketPacketData;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.noData = NoData.decode(reader, reader.uint32());
          break;
        case 3:
          message.createOrderPacket = CreateOrderPacketData.decode(
            reader,
            reader.uint32()
          );
          break;
        case 2:
          message.createPairPacket = CreatePairPacketData.decode(
            reader,
            reader.uint32()
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MarketPacketData {
    const message = { ...baseMarketPacketData } as MarketPacketData;
    if (object.noData !== undefined && object.noData !== null) {
      message.noData = NoData.fromJSON(object.noData);
    } else {
      message.noData = undefined;
    }
    if (
      object.createOrderPacket !== undefined &&
      object.createOrderPacket !== null
    ) {
      message.createOrderPacket = CreateOrderPacketData.fromJSON(
        object.createOrderPacket
      );
    } else {
      message.createOrderPacket = undefined;
    }
    if (
      object.createPairPacket !== undefined &&
      object.createPairPacket !== null
    ) {
      message.createPairPacket = CreatePairPacketData.fromJSON(
        object.createPairPacket
      );
    } else {
      message.createPairPacket = undefined;
    }
    return message;
  },

  toJSON(message: MarketPacketData): unknown {
    const obj: any = {};
    message.noData !== undefined &&
      (obj.noData = message.noData ? NoData.toJSON(message.noData) : undefined);
    message.createOrderPacket !== undefined &&
      (obj.createOrderPacket = message.createOrderPacket
        ? CreateOrderPacketData.toJSON(message.createOrderPacket)
        : undefined);
    message.createPairPacket !== undefined &&
      (obj.createPairPacket = message.createPairPacket
        ? CreatePairPacketData.toJSON(message.createPairPacket)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<MarketPacketData>): MarketPacketData {
    const message = { ...baseMarketPacketData } as MarketPacketData;
    if (object.noData !== undefined && object.noData !== null) {
      message.noData = NoData.fromPartial(object.noData);
    } else {
      message.noData = undefined;
    }
    if (
      object.createOrderPacket !== undefined &&
      object.createOrderPacket !== null
    ) {
      message.createOrderPacket = CreateOrderPacketData.fromPartial(
        object.createOrderPacket
      );
    } else {
      message.createOrderPacket = undefined;
    }
    if (
      object.createPairPacket !== undefined &&
      object.createPairPacket !== null
    ) {
      message.createPairPacket = CreatePairPacketData.fromPartial(
        object.createPairPacket
      );
    } else {
      message.createPairPacket = undefined;
    }
    return message;
  },
};

const baseNoData: object = {};

export const NoData = {
  encode(_: NoData, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): NoData {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseNoData } as NoData;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): NoData {
    const message = { ...baseNoData } as NoData;
    return message;
  },

  toJSON(_: NoData): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<NoData>): NoData {
    const message = { ...baseNoData } as NoData;
    return message;
  },
};

const baseCreateOrderPacketData: object = {
  amountDenom: "",
  amount: 0,
  sourceCoin: "",
  targetCoin: "",
  exchRateDenom: "",
  exchRate: "",
};

export const CreateOrderPacketData = {
  encode(
    message: CreateOrderPacketData,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.amountDenom !== "") {
      writer.uint32(10).string(message.amountDenom);
    }
    if (message.amount !== 0) {
      writer.uint32(16).int32(message.amount);
    }
    if (message.sourceCoin !== "") {
      writer.uint32(26).string(message.sourceCoin);
    }
    if (message.targetCoin !== "") {
      writer.uint32(34).string(message.targetCoin);
    }
    if (message.exchRateDenom !== "") {
      writer.uint32(42).string(message.exchRateDenom);
    }
    if (message.exchRate !== "") {
      writer.uint32(50).string(message.exchRate);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): CreateOrderPacketData {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseCreateOrderPacketData } as CreateOrderPacketData;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.amountDenom = reader.string();
          break;
        case 2:
          message.amount = reader.int32();
          break;
        case 3:
          message.sourceCoin = reader.string();
          break;
        case 4:
          message.targetCoin = reader.string();
          break;
        case 5:
          message.exchRateDenom = reader.string();
          break;
        case 6:
          message.exchRate = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateOrderPacketData {
    const message = { ...baseCreateOrderPacketData } as CreateOrderPacketData;
    if (object.amountDenom !== undefined && object.amountDenom !== null) {
      message.amountDenom = String(object.amountDenom);
    } else {
      message.amountDenom = "";
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = Number(object.amount);
    } else {
      message.amount = 0;
    }
    if (object.sourceCoin !== undefined && object.sourceCoin !== null) {
      message.sourceCoin = String(object.sourceCoin);
    } else {
      message.sourceCoin = "";
    }
    if (object.targetCoin !== undefined && object.targetCoin !== null) {
      message.targetCoin = String(object.targetCoin);
    } else {
      message.targetCoin = "";
    }
    if (object.exchRateDenom !== undefined && object.exchRateDenom !== null) {
      message.exchRateDenom = String(object.exchRateDenom);
    } else {
      message.exchRateDenom = "";
    }
    if (object.exchRate !== undefined && object.exchRate !== null) {
      message.exchRate = String(object.exchRate);
    } else {
      message.exchRate = "";
    }
    return message;
  },

  toJSON(message: CreateOrderPacketData): unknown {
    const obj: any = {};
    message.amountDenom !== undefined &&
      (obj.amountDenom = message.amountDenom);
    message.amount !== undefined && (obj.amount = message.amount);
    message.sourceCoin !== undefined && (obj.sourceCoin = message.sourceCoin);
    message.targetCoin !== undefined && (obj.targetCoin = message.targetCoin);
    message.exchRateDenom !== undefined &&
      (obj.exchRateDenom = message.exchRateDenom);
    message.exchRate !== undefined && (obj.exchRate = message.exchRate);
    return obj;
  },

  fromPartial(
    object: DeepPartial<CreateOrderPacketData>
  ): CreateOrderPacketData {
    const message = { ...baseCreateOrderPacketData } as CreateOrderPacketData;
    if (object.amountDenom !== undefined && object.amountDenom !== null) {
      message.amountDenom = object.amountDenom;
    } else {
      message.amountDenom = "";
    }
    if (object.amount !== undefined && object.amount !== null) {
      message.amount = object.amount;
    } else {
      message.amount = 0;
    }
    if (object.sourceCoin !== undefined && object.sourceCoin !== null) {
      message.sourceCoin = object.sourceCoin;
    } else {
      message.sourceCoin = "";
    }
    if (object.targetCoin !== undefined && object.targetCoin !== null) {
      message.targetCoin = object.targetCoin;
    } else {
      message.targetCoin = "";
    }
    if (object.exchRateDenom !== undefined && object.exchRateDenom !== null) {
      message.exchRateDenom = object.exchRateDenom;
    } else {
      message.exchRateDenom = "";
    }
    if (object.exchRate !== undefined && object.exchRate !== null) {
      message.exchRate = object.exchRate;
    } else {
      message.exchRate = "";
    }
    return message;
  },
};

const baseCreateOrderPacketAck: object = {};

export const CreateOrderPacketAck = {
  encode(_: CreateOrderPacketAck, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): CreateOrderPacketAck {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseCreateOrderPacketAck } as CreateOrderPacketAck;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): CreateOrderPacketAck {
    const message = { ...baseCreateOrderPacketAck } as CreateOrderPacketAck;
    return message;
  },

  toJSON(_: CreateOrderPacketAck): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<CreateOrderPacketAck>): CreateOrderPacketAck {
    const message = { ...baseCreateOrderPacketAck } as CreateOrderPacketAck;
    return message;
  },
};

const baseCreatePairPacketData: object = { sourceDenom: "", targetDenom: "" };

export const CreatePairPacketData = {
  encode(
    message: CreatePairPacketData,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.sourceDenom !== "") {
      writer.uint32(10).string(message.sourceDenom);
    }
    if (message.targetDenom !== "") {
      writer.uint32(18).string(message.targetDenom);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): CreatePairPacketData {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseCreatePairPacketData } as CreatePairPacketData;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sourceDenom = reader.string();
          break;
        case 2:
          message.targetDenom = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreatePairPacketData {
    const message = { ...baseCreatePairPacketData } as CreatePairPacketData;
    if (object.sourceDenom !== undefined && object.sourceDenom !== null) {
      message.sourceDenom = String(object.sourceDenom);
    } else {
      message.sourceDenom = "";
    }
    if (object.targetDenom !== undefined && object.targetDenom !== null) {
      message.targetDenom = String(object.targetDenom);
    } else {
      message.targetDenom = "";
    }
    return message;
  },

  toJSON(message: CreatePairPacketData): unknown {
    const obj: any = {};
    message.sourceDenom !== undefined &&
      (obj.sourceDenom = message.sourceDenom);
    message.targetDenom !== undefined &&
      (obj.targetDenom = message.targetDenom);
    return obj;
  },

  fromPartial(object: DeepPartial<CreatePairPacketData>): CreatePairPacketData {
    const message = { ...baseCreatePairPacketData } as CreatePairPacketData;
    if (object.sourceDenom !== undefined && object.sourceDenom !== null) {
      message.sourceDenom = object.sourceDenom;
    } else {
      message.sourceDenom = "";
    }
    if (object.targetDenom !== undefined && object.targetDenom !== null) {
      message.targetDenom = object.targetDenom;
    } else {
      message.targetDenom = "";
    }
    return message;
  },
};

const baseCreatePairPacketAck: object = {};

export const CreatePairPacketAck = {
  encode(_: CreatePairPacketAck, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): CreatePairPacketAck {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseCreatePairPacketAck } as CreatePairPacketAck;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): CreatePairPacketAck {
    const message = { ...baseCreatePairPacketAck } as CreatePairPacketAck;
    return message;
  },

  toJSON(_: CreatePairPacketAck): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<CreatePairPacketAck>): CreatePairPacketAck {
    const message = { ...baseCreatePairPacketAck } as CreatePairPacketAck;
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
