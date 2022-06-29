import { postReq, getReq } from '../abstract';
import {
  LoginParam,
  SignUpParam,
  UserModifyPasswordParam,
  ClientDetailParam,
  ClientHistoryDealStatsParam,
  ClientProfileParam,
  ClientReviewsParam,
  ProvidersParam,
  SpOwnerProfileParam,
  SpOwnerReviewsParam,
  SpServiceDetailParam,
  PostClientDetailParam,
  PostClientReviewParam,
  ProviderDetailParam,
  ChatHistoryParam,
  clientsParam,
} from '@/api/modules/interface';
import { ProfileResponse } from '@/models/types';

export function postUserSignUp(data: SignUpParam): Promise<any> {
  return postReq({
    url: 'api/v1/userSignUp',
    data: data,
  });
}

export function postUserLogin(data: LoginParam): Promise<any> {
  return postReq({
    url: 'api/v1/userLogin',
    data: data,
  });
}

export function postUserModifyPassword(
  data: UserModifyPasswordParam,
): Promise<any> {
  return postReq({
    url: 'api/v1/user/modifyPassword',
    data: data,
  });
}

export function getClientDetail(data: ClientDetailParam): Promise<any> {
  return getReq({
    url: 'api/v1/clientDetail',
    params: data,
  });
}

export function getClientHistoryDealStats(
  data: ClientHistoryDealStatsParam,
): Promise<any> {
  return getReq({
    url: 'api/v1/clientHistoryDealStats',
    params: data,
  });
}

export function getClientProfile(data: ClientProfileParam): Promise<any> {
  return getReq({
    url: 'api/v1/clientProfile',
    params: data,
  });
}

export function getClientReviews(data: ClientReviewsParam): Promise<any> {
  return getReq({
    url: 'api/v1/clientReviews',
    params: data,
  });
}

export function getProviders(data: ProvidersParam): Promise<any> {
  return getReq({
    url: 'api/v1/providers',
    params: data,
  });
}

export function getClients(data: clientsParam): Promise<any> {
  return getReq({
    url: 'api/v1/clients',
    params: data,
  });
}

export function getSpOwnerProfile(data: SpOwnerProfileParam): Promise<any> {
  return getReq({
    url: 'api/v1/spOwnerProfile',
    params: data,
  });
}

export function getSpOwnerReviews(data: SpOwnerReviewsParam): Promise<any> {
  return getReq({
    url: 'api/v1/spOwnerReviews',
    params: data,
  });
}

export function getSpServiceDetail(data: SpServiceDetailParam): Promise<any> {
  return getReq({
    url: 'api/v1/spServiceDetail',
    params: data,
  });
}

export function postClientDetail(data: PostClientDetailParam): Promise<any> {
  return postReq({
    url: 'api/v1/client/detail',
    data: data,
  });
}

export function postClientReview(data: PostClientReviewParam): Promise<any> {
  return postReq({
    url: 'api/v1/client/review',
    data: data,
  });
}

export function postUserProfile(data: ProfileResponse): Promise<any> {
  return postReq({
    url: 'api/v1/user/profile',
    data: data,
  });
}

export function postProviderDetail(data: ProviderDetailParam): Promise<any> {
  return postReq({
    url: 'api/v1/provider/detail',
    data: data,
  });
}

export function getChatHistory(data: ChatHistoryParam) {
  return getReq({
    url: 'api/v1/chat/history',
    params: data,
  });
}
