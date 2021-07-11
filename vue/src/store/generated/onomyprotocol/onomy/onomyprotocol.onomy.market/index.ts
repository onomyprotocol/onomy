import { txClient, queryClient } from './module'
// @ts-ignore
import { SpVuexError } from '@starport/vuex'

import { DenomTrace } from "./module/types/market/denomTrace"
import { OrderBook } from "./module/types/market/orderBook"
import { MarketPacketData } from "./module/types/market/packet"
import { NoData } from "./module/types/market/packet"
import { CreateOrderPacketData } from "./module/types/market/packet"
import { CreateOrderPacketAck } from "./module/types/market/packet"
import { CreatePairPacketData } from "./module/types/market/packet"
import { CreatePairPacketAck } from "./module/types/market/packet"
import { Order } from "./module/types/market/order"


async function initTxClient(vuexGetters) {
	return await txClient(vuexGetters['common/wallet/signer'], {
		addr: vuexGetters['common/env/apiTendermint']
	})
}

async function initQueryClient(vuexGetters) {
	return await queryClient({
		addr: vuexGetters['common/env/apiCosmos']
	})
}

function getStructure(template) {
	let structure = { fields: [] }
	for (const [key, value] of Object.entries(template)) {
		let field: any = {}
		field.name = key
		field.type = typeof value
		structure.fields.push(field)
	}
	return structure
}

const getDefaultState = () => {
	return {
        DenomTrace: {},
        DenomTraceAll: {},
        OrderBook: {},
        OrderBookAll: {},
        
        _Structure: {
            DenomTrace: getStructure(DenomTrace.fromPartial({})),
            OrderBook: getStructure(OrderBook.fromPartial({})),
            MarketPacketData: getStructure(MarketPacketData.fromPartial({})),
            NoData: getStructure(NoData.fromPartial({})),
            CreateOrderPacketData: getStructure(CreateOrderPacketData.fromPartial({})),
            CreateOrderPacketAck: getStructure(CreateOrderPacketAck.fromPartial({})),
            CreatePairPacketData: getStructure(CreatePairPacketData.fromPartial({})),
            CreatePairPacketAck: getStructure(CreatePairPacketAck.fromPartial({})),
            Order: getStructure(Order.fromPartial({})),
            
		},
		_Subscriptions: new Set(),
	}
}

// initial state
const state = getDefaultState()

export default {
	namespaced: true,
	state,
	mutations: {
		RESET_STATE(state) {
			Object.assign(state, getDefaultState())
		},
		QUERY(state, { query, key, value }) {
			state[query][JSON.stringify(key)] = value
		},
		SUBSCRIBE(state, subscription) {
			state._Subscriptions.add(subscription)
		},
		UNSUBSCRIBE(state, subscription) {
			state._Subscriptions.delete(subscription)
		}
	},
	getters: {
        getDenomTrace: (state) => (params = {}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.DenomTrace[JSON.stringify(params)] ?? {}
		},
        getDenomTraceAll: (state) => (params = {}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.DenomTraceAll[JSON.stringify(params)] ?? {}
		},
        getOrderBook: (state) => (params = {}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.OrderBook[JSON.stringify(params)] ?? {}
		},
        getOrderBookAll: (state) => (params = {}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.OrderBookAll[JSON.stringify(params)] ?? {}
		},
        
		getTypeStructure: (state) => (type) => {
			return state._Structure[type].fields
		}
	},
	actions: {
		init({ dispatch, rootGetters }) {
			console.log('init')
			if (rootGetters['common/env/client']) {
				rootGetters['common/env/client'].on('newblock', () => {
					dispatch('StoreUpdate')
				})
			}
		},
		resetState({ commit }) {
			commit('RESET_STATE')
		},
		unsubscribe({ commit }, subscription) {
			commit('UNSUBSCRIBE', subscription)
		},
		async StoreUpdate({ state, dispatch }) {
			state._Subscriptions.forEach((subscription) => {
				dispatch(subscription.action, subscription.payload)
			})
		},
		async QueryDenomTrace({ commit, rootGetters, getters }, { options: { subscribe = false , all = false}, params: {...key}, query=null }) {
			try {
				
				let value = query?(await (await initQueryClient(rootGetters)).queryDenomTrace( key.index,  query)).data:(await (await initQueryClient(rootGetters)).queryDenomTrace( key.index )).data
				
				commit('QUERY', { query: 'DenomTrace', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDenomTrace', payload: { options: { all }, params: {...key},query }})
				return getters['getDenomTrace']( { params: {...key}, query}) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDenomTrace', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryDenomTraceAll({ commit, rootGetters, getters }, { options: { subscribe = false , all = false}, params: {...key}, query=null }) {
			try {
				
				let value = query?(await (await initQueryClient(rootGetters)).queryDenomTraceAll( query)).data:(await (await initQueryClient(rootGetters)).queryDenomTraceAll()).data
				
				while (all && (<any> value).pagination && (<any> value).pagination.nextKey!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryDenomTraceAll({...query, 'pagination.key':(<any> value).pagination.nextKey})).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				
				commit('QUERY', { query: 'DenomTraceAll', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryDenomTraceAll', payload: { options: { all }, params: {...key},query }})
				return getters['getDenomTraceAll']( { params: {...key}, query}) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryDenomTraceAll', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryOrderBook({ commit, rootGetters, getters }, { options: { subscribe = false , all = false}, params: {...key}, query=null }) {
			try {
				
				let value = query?(await (await initQueryClient(rootGetters)).queryOrderBook( key.index,  query)).data:(await (await initQueryClient(rootGetters)).queryOrderBook( key.index )).data
				
				commit('QUERY', { query: 'OrderBook', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryOrderBook', payload: { options: { all }, params: {...key},query }})
				return getters['getOrderBook']( { params: {...key}, query}) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryOrderBook', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		async QueryOrderBookAll({ commit, rootGetters, getters }, { options: { subscribe = false , all = false}, params: {...key}, query=null }) {
			try {
				
				let value = query?(await (await initQueryClient(rootGetters)).queryOrderBookAll( query)).data:(await (await initQueryClient(rootGetters)).queryOrderBookAll()).data
				
				while (all && (<any> value).pagination && (<any> value).pagination.nextKey!=null) {
					let next_values=(await (await initQueryClient(rootGetters)).queryOrderBookAll({...query, 'pagination.key':(<any> value).pagination.nextKey})).data
					for (let prop of Object.keys(next_values)) {
						if (Array.isArray(next_values[prop])) {
							value[prop]=[...value[prop], ...next_values[prop]]
						}else{
							value[prop]=next_values[prop]
						}
					}
				}
				
				commit('QUERY', { query: 'OrderBookAll', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryOrderBookAll', payload: { options: { all }, params: {...key},query }})
				return getters['getOrderBookAll']( { params: {...key}, query}) ?? {}
			} catch (e) {
				console.error(new SpVuexError('QueryClient:QueryOrderBookAll', 'API Node Unavailable. Could not perform query.'))
				return {}
			}
		},
		
		async sendMsgSendCreatePair({ rootGetters }, { value, fee, memo }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgSendCreatePair(value)
				const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], {fee: { amount: fee, 
  gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgSendCreatePair:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgSendCreatePair:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgCancelOrder({ rootGetters }, { value, fee, memo }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgCancelOrder(value)
				const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], {fee: { amount: fee, 
  gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgCancelOrder:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgCancelOrder:Send', 'Could not broadcast Tx.')
				}
			}
		},
		async sendMsgSendCreateOrder({ rootGetters }, { value, fee, memo }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgSendCreateOrder(value)
				const result = await (await initTxClient(rootGetters)).signAndBroadcast([msg], {fee: { amount: fee, 
  gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgSendCreateOrder:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgSendCreateOrder:Send', 'Could not broadcast Tx.')
				}
			}
		},
		
		async MsgSendCreatePair({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgSendCreatePair(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgSendCreatePair:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgSendCreatePair:Create', 'Could not create message.')
				}
			}
		},
		async MsgCancelOrder({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgCancelOrder(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgCancelOrder:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgCancelOrder:Create', 'Could not create message.')
				}
			}
		},
		async MsgSendCreateOrder({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).msgSendCreateOrder(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:MsgSendCreateOrder:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgSendCreateOrder:Create', 'Could not create message.')
				}
			}
		},
		
	}
}
