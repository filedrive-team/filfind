export interface SignUpParam {
  address: string;
  email: string;
  message: string;
  name: string;
  password: string;
  signature: string;
  type: string;
}

export interface HeadersParam {
  Authorization: string;
}

export interface LoginParam {
  email: string;
  password: string;
}

export interface UserModifyPasswordParam {
  new_password: string;
  password: string;
}

export interface ClientDetailParam {
  address_id: string;
}

export interface ClientHistoryDealStatsParam {
  page: number;
  page_size: number;
  address_id: string;
}

export interface ClientProfileParam {
  address_id: string;
}

export interface ClientReviewsParam {
  page: number;
  page_size: number;
  address_id: string;
}

export interface ProvidersParam {
  page: number;
  page_size: number;
  sort_by: string;
  order: string;
  sps_status: string;
  region: string;
  raw_power_range: string;
  storage_success_rate_range: string;
  reputation_score_range: string;
  review_score_range: string;
  search: string;
}

export interface SpOwnerProfileParam {
  address_id: string;
}

export interface SpOwnerReviewsParam {
  page: number;
  page_size: number;
  address_id: string;
}

export interface SpServiceDetailParam {
  address_id: string;
}

export interface PostClientDetailParam {
  bandwidth: string;
  monthly_storage: string;
  service_requirement: string;
  use_case: string;
}

export interface PostClientReviewParam {
  title: string;
  content: string;
  provider: string;
  score: number;
}

export interface ProviderDetailParam {
  address: string;
  available_deals: string;
  bandwidth: string;
  certification: string;
  experience: string;
  is_member: string;
  parallel_deals: string;
  renewable_energy: string;
  sealing_speed: string;
}

export interface ChatHistoryParam {
  partner: string | undefined;
  before: number | null;
  limit: number;
}

export interface clientsParam {
  page: number;
  page_size: number;
  sort_by: string;
  order: string;
  search: string;
}
