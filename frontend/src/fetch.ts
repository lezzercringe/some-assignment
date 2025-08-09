import config from "./config";

export type Result<T, E> = { ok: true; value: T } | { ok: false; error: E };

export async function FetchOrder(id: string): Promise<Result<Order, FetchError>> {
    console.log(config.apiUrl)
    const url = new URL(`/order/${id}`, config.apiUrl).toString();

    try {
        const response = await fetch(url);

        if (!response.ok) {
            try {
                const json = await response.json();
                return { ok: false, error: { type: "api", message: (json as ApiError).message } };
            } catch {
                return { ok: false, error: { type: "http", code: response.status, statusText: response.statusText } };
            }
        }

        const data = await response.json();
        return { ok: true, value: data as Order };
    } catch (error) {
        return { ok: false, error: { type: "network", error: error as Error } };
    }
}

export type FetchError = HttpError | NetworkError | ApiError;

type HttpError = {
    type: "http";
    code: number;
    statusText: string;
};

type NetworkError = {
    type: "network";
    error: Error;
};


type ApiError = {
    type: "api";
    message: string;
};

export type Order = {
    order_uid: string;
    track_number: string;
    entry: string;
    delivery: {
        name: string;
        phone: string;
        zip: string;
        city: string;
        address: string;
        region: string;
        email: string;
    };
    payment: {
        transaction: string;
        request_id: string;
        currency: string;
        provider: string;
        amount: number;
        payment_dt: number;
        bank: string;
        delivery_cost: number;
        goods_total: number;
        custom_fee: number;
    };
    items: {
        chrt_id: number;
        track_number: string;
        price: number;
        rid: string;
        name: string;
        sale: number;
        size: string;
        total_price: number;
        nm_id: number;
        brand: string;
        status: number;
    }[];
    locale: string;
    internal_signature: string;
    customer_id: string;
    delivery_service: string;
    shardkey: string;
    sm_id: number;
    date_created: string;
    oof_shard: string;
};
