/* eslint-disable */
import { Reader, Writer } from "protobufjs/minimal";
import { DenomTrace } from "../market/denomTrace";
import {
  PageRequest,
  PageResponse,
} from "../cosmos/base/query/v1beta1/pagination";
import { OrderBook } from "../market/orderBook";

export const protobufPackage = "onomyprotocol.onomy.market";

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

const baseQueryGetDenomTraceRequest: object = { index: "" };

export const QueryGetDenomTraceRequest = {
  encode(
    message: QueryGetDenomTraceRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.index !== "") {
      writer.uint32(10).string(message.index);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetDenomTraceRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetDenomTraceRequest,
    } as QueryGetDenomTraceRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.index = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetDenomTraceRequest {
    const message = {
      ...baseQueryGetDenomTraceRequest,
    } as QueryGetDenomTraceRequest;
    if (object.index !== undefined && object.index !== null) {
      message.index = String(object.index);
    } else {
      message.index = "";
    }
    return message;
  },

  toJSON(message: QueryGetDenomTraceRequest): unknown {
    const obj: any = {};
    message.index !== undefined && (obj.index = message.index);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetDenomTraceRequest>
  ): QueryGetDenomTraceRequest {
    const message = {
      ...baseQueryGetDenomTraceRequest,
    } as QueryGetDenomTraceRequest;
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    } else {
      message.index = "";
    }
    return message;
  },
};

const baseQueryGetDenomTraceResponse: object = {};

export const QueryGetDenomTraceResponse = {
  encode(
    message: QueryGetDenomTraceResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.DenomTrace !== undefined) {
      DenomTrace.encode(message.DenomTrace, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetDenomTraceResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetDenomTraceResponse,
    } as QueryGetDenomTraceResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.DenomTrace = DenomTrace.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetDenomTraceResponse {
    const message = {
      ...baseQueryGetDenomTraceResponse,
    } as QueryGetDenomTraceResponse;
    if (object.DenomTrace !== undefined && object.DenomTrace !== null) {
      message.DenomTrace = DenomTrace.fromJSON(object.DenomTrace);
    } else {
      message.DenomTrace = undefined;
    }
    return message;
  },

  toJSON(message: QueryGetDenomTraceResponse): unknown {
    const obj: any = {};
    message.DenomTrace !== undefined &&
      (obj.DenomTrace = message.DenomTrace
        ? DenomTrace.toJSON(message.DenomTrace)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetDenomTraceResponse>
  ): QueryGetDenomTraceResponse {
    const message = {
      ...baseQueryGetDenomTraceResponse,
    } as QueryGetDenomTraceResponse;
    if (object.DenomTrace !== undefined && object.DenomTrace !== null) {
      message.DenomTrace = DenomTrace.fromPartial(object.DenomTrace);
    } else {
      message.DenomTrace = undefined;
    }
    return message;
  },
};

const baseQueryAllDenomTraceRequest: object = {};

export const QueryAllDenomTraceRequest = {
  encode(
    message: QueryAllDenomTraceRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllDenomTraceRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllDenomTraceRequest,
    } as QueryAllDenomTraceRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllDenomTraceRequest {
    const message = {
      ...baseQueryAllDenomTraceRequest,
    } as QueryAllDenomTraceRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllDenomTraceRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageRequest.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllDenomTraceRequest>
  ): QueryAllDenomTraceRequest {
    const message = {
      ...baseQueryAllDenomTraceRequest,
    } as QueryAllDenomTraceRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryAllDenomTraceResponse: object = {};

export const QueryAllDenomTraceResponse = {
  encode(
    message: QueryAllDenomTraceResponse,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.DenomTrace) {
      DenomTrace.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllDenomTraceResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllDenomTraceResponse,
    } as QueryAllDenomTraceResponse;
    message.DenomTrace = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.DenomTrace.push(DenomTrace.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllDenomTraceResponse {
    const message = {
      ...baseQueryAllDenomTraceResponse,
    } as QueryAllDenomTraceResponse;
    message.DenomTrace = [];
    if (object.DenomTrace !== undefined && object.DenomTrace !== null) {
      for (const e of object.DenomTrace) {
        message.DenomTrace.push(DenomTrace.fromJSON(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllDenomTraceResponse): unknown {
    const obj: any = {};
    if (message.DenomTrace) {
      obj.DenomTrace = message.DenomTrace.map((e) =>
        e ? DenomTrace.toJSON(e) : undefined
      );
    } else {
      obj.DenomTrace = [];
    }
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllDenomTraceResponse>
  ): QueryAllDenomTraceResponse {
    const message = {
      ...baseQueryAllDenomTraceResponse,
    } as QueryAllDenomTraceResponse;
    message.DenomTrace = [];
    if (object.DenomTrace !== undefined && object.DenomTrace !== null) {
      for (const e of object.DenomTrace) {
        message.DenomTrace.push(DenomTrace.fromPartial(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryGetOrderBookRequest: object = { index: "" };

export const QueryGetOrderBookRequest = {
  encode(
    message: QueryGetOrderBookRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.index !== "") {
      writer.uint32(10).string(message.index);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetOrderBookRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetOrderBookRequest,
    } as QueryGetOrderBookRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.index = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetOrderBookRequest {
    const message = {
      ...baseQueryGetOrderBookRequest,
    } as QueryGetOrderBookRequest;
    if (object.index !== undefined && object.index !== null) {
      message.index = String(object.index);
    } else {
      message.index = "";
    }
    return message;
  },

  toJSON(message: QueryGetOrderBookRequest): unknown {
    const obj: any = {};
    message.index !== undefined && (obj.index = message.index);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetOrderBookRequest>
  ): QueryGetOrderBookRequest {
    const message = {
      ...baseQueryGetOrderBookRequest,
    } as QueryGetOrderBookRequest;
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    } else {
      message.index = "";
    }
    return message;
  },
};

const baseQueryGetOrderBookResponse: object = {};

export const QueryGetOrderBookResponse = {
  encode(
    message: QueryGetOrderBookResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.OrderBook !== undefined) {
      OrderBook.encode(message.OrderBook, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryGetOrderBookResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryGetOrderBookResponse,
    } as QueryGetOrderBookResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.OrderBook = OrderBook.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetOrderBookResponse {
    const message = {
      ...baseQueryGetOrderBookResponse,
    } as QueryGetOrderBookResponse;
    if (object.OrderBook !== undefined && object.OrderBook !== null) {
      message.OrderBook = OrderBook.fromJSON(object.OrderBook);
    } else {
      message.OrderBook = undefined;
    }
    return message;
  },

  toJSON(message: QueryGetOrderBookResponse): unknown {
    const obj: any = {};
    message.OrderBook !== undefined &&
      (obj.OrderBook = message.OrderBook
        ? OrderBook.toJSON(message.OrderBook)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryGetOrderBookResponse>
  ): QueryGetOrderBookResponse {
    const message = {
      ...baseQueryGetOrderBookResponse,
    } as QueryGetOrderBookResponse;
    if (object.OrderBook !== undefined && object.OrderBook !== null) {
      message.OrderBook = OrderBook.fromPartial(object.OrderBook);
    } else {
      message.OrderBook = undefined;
    }
    return message;
  },
};

const baseQueryAllOrderBookRequest: object = {};

export const QueryAllOrderBookRequest = {
  encode(
    message: QueryAllOrderBookRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllOrderBookRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllOrderBookRequest,
    } as QueryAllOrderBookRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllOrderBookRequest {
    const message = {
      ...baseQueryAllOrderBookRequest,
    } as QueryAllOrderBookRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllOrderBookRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageRequest.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllOrderBookRequest>
  ): QueryAllOrderBookRequest {
    const message = {
      ...baseQueryAllOrderBookRequest,
    } as QueryAllOrderBookRequest;
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

const baseQueryAllOrderBookResponse: object = {};

export const QueryAllOrderBookResponse = {
  encode(
    message: QueryAllOrderBookResponse,
    writer: Writer = Writer.create()
  ): Writer {
    for (const v of message.OrderBook) {
      OrderBook.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(
        message.pagination,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): QueryAllOrderBookResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseQueryAllOrderBookResponse,
    } as QueryAllOrderBookResponse;
    message.OrderBook = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.OrderBook.push(OrderBook.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllOrderBookResponse {
    const message = {
      ...baseQueryAllOrderBookResponse,
    } as QueryAllOrderBookResponse;
    message.OrderBook = [];
    if (object.OrderBook !== undefined && object.OrderBook !== null) {
      for (const e of object.OrderBook) {
        message.OrderBook.push(OrderBook.fromJSON(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromJSON(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },

  toJSON(message: QueryAllOrderBookResponse): unknown {
    const obj: any = {};
    if (message.OrderBook) {
      obj.OrderBook = message.OrderBook.map((e) =>
        e ? OrderBook.toJSON(e) : undefined
      );
    } else {
      obj.OrderBook = [];
    }
    message.pagination !== undefined &&
      (obj.pagination = message.pagination
        ? PageResponse.toJSON(message.pagination)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<QueryAllOrderBookResponse>
  ): QueryAllOrderBookResponse {
    const message = {
      ...baseQueryAllOrderBookResponse,
    } as QueryAllOrderBookResponse;
    message.OrderBook = [];
    if (object.OrderBook !== undefined && object.OrderBook !== null) {
      for (const e of object.OrderBook) {
        message.OrderBook.push(OrderBook.fromPartial(e));
      }
    }
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromPartial(object.pagination);
    } else {
      message.pagination = undefined;
    }
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** this line is used by starport scaffolding # 2 */
  DenomTrace(
    request: QueryGetDenomTraceRequest
  ): Promise<QueryGetDenomTraceResponse>;
  DenomTraceAll(
    request: QueryAllDenomTraceRequest
  ): Promise<QueryAllDenomTraceResponse>;
  OrderBook(
    request: QueryGetOrderBookRequest
  ): Promise<QueryGetOrderBookResponse>;
  OrderBookAll(
    request: QueryAllOrderBookRequest
  ): Promise<QueryAllOrderBookResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  DenomTrace(
    request: QueryGetDenomTraceRequest
  ): Promise<QueryGetDenomTraceResponse> {
    const data = QueryGetDenomTraceRequest.encode(request).finish();
    const promise = this.rpc.request(
      "onomyprotocol.onomy.market.Query",
      "DenomTrace",
      data
    );
    return promise.then((data) =>
      QueryGetDenomTraceResponse.decode(new Reader(data))
    );
  }

  DenomTraceAll(
    request: QueryAllDenomTraceRequest
  ): Promise<QueryAllDenomTraceResponse> {
    const data = QueryAllDenomTraceRequest.encode(request).finish();
    const promise = this.rpc.request(
      "onomyprotocol.onomy.market.Query",
      "DenomTraceAll",
      data
    );
    return promise.then((data) =>
      QueryAllDenomTraceResponse.decode(new Reader(data))
    );
  }

  OrderBook(
    request: QueryGetOrderBookRequest
  ): Promise<QueryGetOrderBookResponse> {
    const data = QueryGetOrderBookRequest.encode(request).finish();
    const promise = this.rpc.request(
      "onomyprotocol.onomy.market.Query",
      "OrderBook",
      data
    );
    return promise.then((data) =>
      QueryGetOrderBookResponse.decode(new Reader(data))
    );
  }

  OrderBookAll(
    request: QueryAllOrderBookRequest
  ): Promise<QueryAllOrderBookResponse> {
    const data = QueryAllOrderBookRequest.encode(request).finish();
    const promise = this.rpc.request(
      "onomyprotocol.onomy.market.Query",
      "OrderBookAll",
      data
    );
    return promise.then((data) =>
      QueryAllOrderBookResponse.decode(new Reader(data))
    );
  }
}

interface Rpc {
  request(
    service: string,
    method: string,
    data: Uint8Array
  ): Promise<Uint8Array>;
}

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
