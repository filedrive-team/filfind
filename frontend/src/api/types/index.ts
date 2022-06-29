export type Method = 'GET' | 'POST' | 'PUT' | 'DELETE';
export type ResponseType =
  | 'arraybuffer'
  | 'blob'
  | 'document'
  | 'json'
  | 'text'
  | 'stream';

export interface AxiosRequest {
  baseURL?: string;
  url: string | undefined;
  data?: object;
  params?: object;
  method?: Method;
  headers?: {};
  timeout?: number;
  responseType?: ResponseType;
}

export interface AxiosResponse {
  data: any;
  headers: object;
  request?: object;
  status: number;
  statusText: string;
  config: AxiosRequest;
}

export interface CustomResponse {
  readonly status: boolean;
  readonly message: string;
  data: any;
  origin?: any;
}

export interface Dashboard {
  TVL: number;
  cumulativeTVL: number;
  cumulativeTraders: number;
  cumulativeTransactions: number;
  numberOfTraders: number;
  numberOfTransactions: number;
  price: number;
  tradingPairs: number;
  tradingVolume: number;
  snapshotTime: string;
}

export interface spServiceDetailResponse {
  address: string;
  available_deals: string;
  bandwidth: string;
  certification: string;
  experience: string;
  first_deal_time: number;
  is_member: string;
  iso_code: string;
  max_piece_size: string;
  min_piece_size: string;
  name: string;
  owner: string;
  parallel_deals: string;
  price: string;
  quality_adj_power: string;
  raw_power: string;
  region: string;
  renewable_energy: string;
  reputation_score: number;
  retrieval_success_rate: number;
  review_score: number;
  reviews: number;
  sealing_speed: string;
  storage_deals: number;
  storage_success_rate: number;
  verified_price: string;
}
