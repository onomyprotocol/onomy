/* eslint-disable */
import { Reader, Writer } from "protobufjs/minimal";
import { DenomTrace } from "../market/denomTrace";
import { PageRequest, PageResponse, } from "../cosmos/base/query/v1beta1/pagination";
import { OrderBook } from "../market/orderBook";
export const protobufPackage = "onomyprotocol.onomy.market";
const baseQueryGetDenomTraceRequest = { index: "" };
export const QueryGetDenomTraceRequest = {
    encode(message, writer = Writer.create()) {
        if (message.index !== "") {
            writer.uint32(10).string(message.index);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryGetDenomTraceRequest,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryGetDenomTraceRequest,
        };
        if (object.index !== undefined && object.index !== null) {
            message.index = String(object.index);
        }
        else {
            message.index = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.index !== undefined && (obj.index = message.index);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryGetDenomTraceRequest,
        };
        if (object.index !== undefined && object.index !== null) {
            message.index = object.index;
        }
        else {
            message.index = "";
        }
        return message;
    },
};
const baseQueryGetDenomTraceResponse = {};
export const QueryGetDenomTraceResponse = {
    encode(message, writer = Writer.create()) {
        if (message.DenomTrace !== undefined) {
            DenomTrace.encode(message.DenomTrace, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryGetDenomTraceResponse,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryGetDenomTraceResponse,
        };
        if (object.DenomTrace !== undefined && object.DenomTrace !== null) {
            message.DenomTrace = DenomTrace.fromJSON(object.DenomTrace);
        }
        else {
            message.DenomTrace = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.DenomTrace !== undefined &&
            (obj.DenomTrace = message.DenomTrace
                ? DenomTrace.toJSON(message.DenomTrace)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryGetDenomTraceResponse,
        };
        if (object.DenomTrace !== undefined && object.DenomTrace !== null) {
            message.DenomTrace = DenomTrace.fromPartial(object.DenomTrace);
        }
        else {
            message.DenomTrace = undefined;
        }
        return message;
    },
};
const baseQueryAllDenomTraceRequest = {};
export const QueryAllDenomTraceRequest = {
    encode(message, writer = Writer.create()) {
        if (message.pagination !== undefined) {
            PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryAllDenomTraceRequest,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryAllDenomTraceRequest,
        };
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageRequest.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? PageRequest.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryAllDenomTraceRequest,
        };
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageRequest.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
};
const baseQueryAllDenomTraceResponse = {};
export const QueryAllDenomTraceResponse = {
    encode(message, writer = Writer.create()) {
        for (const v of message.DenomTrace) {
            DenomTrace.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryAllDenomTraceResponse,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryAllDenomTraceResponse,
        };
        message.DenomTrace = [];
        if (object.DenomTrace !== undefined && object.DenomTrace !== null) {
            for (const e of object.DenomTrace) {
                message.DenomTrace.push(DenomTrace.fromJSON(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageResponse.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        if (message.DenomTrace) {
            obj.DenomTrace = message.DenomTrace.map((e) => e ? DenomTrace.toJSON(e) : undefined);
        }
        else {
            obj.DenomTrace = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? PageResponse.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryAllDenomTraceResponse,
        };
        message.DenomTrace = [];
        if (object.DenomTrace !== undefined && object.DenomTrace !== null) {
            for (const e of object.DenomTrace) {
                message.DenomTrace.push(DenomTrace.fromPartial(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageResponse.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
};
const baseQueryGetOrderBookRequest = { index: "" };
export const QueryGetOrderBookRequest = {
    encode(message, writer = Writer.create()) {
        if (message.index !== "") {
            writer.uint32(10).string(message.index);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryGetOrderBookRequest,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryGetOrderBookRequest,
        };
        if (object.index !== undefined && object.index !== null) {
            message.index = String(object.index);
        }
        else {
            message.index = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.index !== undefined && (obj.index = message.index);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryGetOrderBookRequest,
        };
        if (object.index !== undefined && object.index !== null) {
            message.index = object.index;
        }
        else {
            message.index = "";
        }
        return message;
    },
};
const baseQueryGetOrderBookResponse = {};
export const QueryGetOrderBookResponse = {
    encode(message, writer = Writer.create()) {
        if (message.OrderBook !== undefined) {
            OrderBook.encode(message.OrderBook, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryGetOrderBookResponse,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryGetOrderBookResponse,
        };
        if (object.OrderBook !== undefined && object.OrderBook !== null) {
            message.OrderBook = OrderBook.fromJSON(object.OrderBook);
        }
        else {
            message.OrderBook = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.OrderBook !== undefined &&
            (obj.OrderBook = message.OrderBook
                ? OrderBook.toJSON(message.OrderBook)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryGetOrderBookResponse,
        };
        if (object.OrderBook !== undefined && object.OrderBook !== null) {
            message.OrderBook = OrderBook.fromPartial(object.OrderBook);
        }
        else {
            message.OrderBook = undefined;
        }
        return message;
    },
};
const baseQueryAllOrderBookRequest = {};
export const QueryAllOrderBookRequest = {
    encode(message, writer = Writer.create()) {
        if (message.pagination !== undefined) {
            PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryAllOrderBookRequest,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryAllOrderBookRequest,
        };
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageRequest.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? PageRequest.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryAllOrderBookRequest,
        };
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageRequest.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
};
const baseQueryAllOrderBookResponse = {};
export const QueryAllOrderBookResponse = {
    encode(message, writer = Writer.create()) {
        for (const v of message.OrderBook) {
            OrderBook.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryAllOrderBookResponse,
        };
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
    fromJSON(object) {
        const message = {
            ...baseQueryAllOrderBookResponse,
        };
        message.OrderBook = [];
        if (object.OrderBook !== undefined && object.OrderBook !== null) {
            for (const e of object.OrderBook) {
                message.OrderBook.push(OrderBook.fromJSON(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageResponse.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        if (message.OrderBook) {
            obj.OrderBook = message.OrderBook.map((e) => e ? OrderBook.toJSON(e) : undefined);
        }
        else {
            obj.OrderBook = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? PageResponse.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryAllOrderBookResponse,
        };
        message.OrderBook = [];
        if (object.OrderBook !== undefined && object.OrderBook !== null) {
            for (const e of object.OrderBook) {
                message.OrderBook.push(OrderBook.fromPartial(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageResponse.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
};
export class QueryClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
    }
    DenomTrace(request) {
        const data = QueryGetDenomTraceRequest.encode(request).finish();
        const promise = this.rpc.request("onomyprotocol.onomy.market.Query", "DenomTrace", data);
        return promise.then((data) => QueryGetDenomTraceResponse.decode(new Reader(data)));
    }
    DenomTraceAll(request) {
        const data = QueryAllDenomTraceRequest.encode(request).finish();
        const promise = this.rpc.request("onomyprotocol.onomy.market.Query", "DenomTraceAll", data);
        return promise.then((data) => QueryAllDenomTraceResponse.decode(new Reader(data)));
    }
    OrderBook(request) {
        const data = QueryGetOrderBookRequest.encode(request).finish();
        const promise = this.rpc.request("onomyprotocol.onomy.market.Query", "OrderBook", data);
        return promise.then((data) => QueryGetOrderBookResponse.decode(new Reader(data)));
    }
    OrderBookAll(request) {
        const data = QueryAllOrderBookRequest.encode(request).finish();
        const promise = this.rpc.request("onomyprotocol.onomy.market.Query", "OrderBookAll", data);
        return promise.then((data) => QueryAllOrderBookResponse.decode(new Reader(data)));
    }
}
