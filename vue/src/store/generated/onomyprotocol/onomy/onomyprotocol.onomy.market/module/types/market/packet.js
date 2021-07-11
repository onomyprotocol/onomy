/* eslint-disable */
import { Writer, Reader } from "protobufjs/minimal";
export const protobufPackage = "onomyprotocol.onomy.market";
const baseMarketPacketData = {};
export const MarketPacketData = {
    encode(message, writer = Writer.create()) {
        if (message.noData !== undefined) {
            NoData.encode(message.noData, writer.uint32(10).fork()).ldelim();
        }
        if (message.createOrderPacket !== undefined) {
            CreateOrderPacketData.encode(message.createOrderPacket, writer.uint32(26).fork()).ldelim();
        }
        if (message.createPairPacket !== undefined) {
            CreatePairPacketData.encode(message.createPairPacket, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMarketPacketData };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.noData = NoData.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.createOrderPacket = CreateOrderPacketData.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.createPairPacket = CreatePairPacketData.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseMarketPacketData };
        if (object.noData !== undefined && object.noData !== null) {
            message.noData = NoData.fromJSON(object.noData);
        }
        else {
            message.noData = undefined;
        }
        if (object.createOrderPacket !== undefined &&
            object.createOrderPacket !== null) {
            message.createOrderPacket = CreateOrderPacketData.fromJSON(object.createOrderPacket);
        }
        else {
            message.createOrderPacket = undefined;
        }
        if (object.createPairPacket !== undefined &&
            object.createPairPacket !== null) {
            message.createPairPacket = CreatePairPacketData.fromJSON(object.createPairPacket);
        }
        else {
            message.createPairPacket = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
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
    fromPartial(object) {
        const message = { ...baseMarketPacketData };
        if (object.noData !== undefined && object.noData !== null) {
            message.noData = NoData.fromPartial(object.noData);
        }
        else {
            message.noData = undefined;
        }
        if (object.createOrderPacket !== undefined &&
            object.createOrderPacket !== null) {
            message.createOrderPacket = CreateOrderPacketData.fromPartial(object.createOrderPacket);
        }
        else {
            message.createOrderPacket = undefined;
        }
        if (object.createPairPacket !== undefined &&
            object.createPairPacket !== null) {
            message.createPairPacket = CreatePairPacketData.fromPartial(object.createPairPacket);
        }
        else {
            message.createPairPacket = undefined;
        }
        return message;
    },
};
const baseNoData = {};
export const NoData = {
    encode(_, writer = Writer.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseNoData };
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
    fromJSON(_) {
        const message = { ...baseNoData };
        return message;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = { ...baseNoData };
        return message;
    },
};
const baseCreateOrderPacketData = {
    amountDenom: "",
    amount: 0,
    sourceCoin: "",
    targetCoin: "",
    exchRateDenom: "",
    exchRate: "",
};
export const CreateOrderPacketData = {
    encode(message, writer = Writer.create()) {
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
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseCreateOrderPacketData };
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
    fromJSON(object) {
        const message = { ...baseCreateOrderPacketData };
        if (object.amountDenom !== undefined && object.amountDenom !== null) {
            message.amountDenom = String(object.amountDenom);
        }
        else {
            message.amountDenom = "";
        }
        if (object.amount !== undefined && object.amount !== null) {
            message.amount = Number(object.amount);
        }
        else {
            message.amount = 0;
        }
        if (object.sourceCoin !== undefined && object.sourceCoin !== null) {
            message.sourceCoin = String(object.sourceCoin);
        }
        else {
            message.sourceCoin = "";
        }
        if (object.targetCoin !== undefined && object.targetCoin !== null) {
            message.targetCoin = String(object.targetCoin);
        }
        else {
            message.targetCoin = "";
        }
        if (object.exchRateDenom !== undefined && object.exchRateDenom !== null) {
            message.exchRateDenom = String(object.exchRateDenom);
        }
        else {
            message.exchRateDenom = "";
        }
        if (object.exchRate !== undefined && object.exchRate !== null) {
            message.exchRate = String(object.exchRate);
        }
        else {
            message.exchRate = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
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
    fromPartial(object) {
        const message = { ...baseCreateOrderPacketData };
        if (object.amountDenom !== undefined && object.amountDenom !== null) {
            message.amountDenom = object.amountDenom;
        }
        else {
            message.amountDenom = "";
        }
        if (object.amount !== undefined && object.amount !== null) {
            message.amount = object.amount;
        }
        else {
            message.amount = 0;
        }
        if (object.sourceCoin !== undefined && object.sourceCoin !== null) {
            message.sourceCoin = object.sourceCoin;
        }
        else {
            message.sourceCoin = "";
        }
        if (object.targetCoin !== undefined && object.targetCoin !== null) {
            message.targetCoin = object.targetCoin;
        }
        else {
            message.targetCoin = "";
        }
        if (object.exchRateDenom !== undefined && object.exchRateDenom !== null) {
            message.exchRateDenom = object.exchRateDenom;
        }
        else {
            message.exchRateDenom = "";
        }
        if (object.exchRate !== undefined && object.exchRate !== null) {
            message.exchRate = object.exchRate;
        }
        else {
            message.exchRate = "";
        }
        return message;
    },
};
const baseCreateOrderPacketAck = {};
export const CreateOrderPacketAck = {
    encode(_, writer = Writer.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseCreateOrderPacketAck };
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
    fromJSON(_) {
        const message = { ...baseCreateOrderPacketAck };
        return message;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = { ...baseCreateOrderPacketAck };
        return message;
    },
};
const baseCreatePairPacketData = { sourceDenom: "", targetDenom: "" };
export const CreatePairPacketData = {
    encode(message, writer = Writer.create()) {
        if (message.sourceDenom !== "") {
            writer.uint32(10).string(message.sourceDenom);
        }
        if (message.targetDenom !== "") {
            writer.uint32(18).string(message.targetDenom);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseCreatePairPacketData };
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
    fromJSON(object) {
        const message = { ...baseCreatePairPacketData };
        if (object.sourceDenom !== undefined && object.sourceDenom !== null) {
            message.sourceDenom = String(object.sourceDenom);
        }
        else {
            message.sourceDenom = "";
        }
        if (object.targetDenom !== undefined && object.targetDenom !== null) {
            message.targetDenom = String(object.targetDenom);
        }
        else {
            message.targetDenom = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.sourceDenom !== undefined &&
            (obj.sourceDenom = message.sourceDenom);
        message.targetDenom !== undefined &&
            (obj.targetDenom = message.targetDenom);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseCreatePairPacketData };
        if (object.sourceDenom !== undefined && object.sourceDenom !== null) {
            message.sourceDenom = object.sourceDenom;
        }
        else {
            message.sourceDenom = "";
        }
        if (object.targetDenom !== undefined && object.targetDenom !== null) {
            message.targetDenom = object.targetDenom;
        }
        else {
            message.targetDenom = "";
        }
        return message;
    },
};
const baseCreatePairPacketAck = {};
export const CreatePairPacketAck = {
    encode(_, writer = Writer.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseCreatePairPacketAck };
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
    fromJSON(_) {
        const message = { ...baseCreatePairPacketAck };
        return message;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = { ...baseCreatePairPacketAck };
        return message;
    },
};
