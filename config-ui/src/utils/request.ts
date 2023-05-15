/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import type { AxiosRequestConfig } from 'axios';
import axios from 'axios';
import { history } from '@/utils/history';

import { DEVLAKE_ENDPOINT } from '@/config';
import { toast } from '@/components/toast';

const instance = axios.create({
  baseURL: DEVLAKE_ENDPOINT,
});

var refreshingToken: Promise<any> | null = null;

instance.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error.response?.status;

    if (status === 401) {
      toast.error('Please login first');
      history.push('/login');
    }

    if (status === 403) {
      var refreshToken = localStorage.getItem('refreshToken');
      if (refreshToken) {
        refreshingToken =
          refreshingToken ||
          request('/login/refreshtoken', {
            method: 'POST',
            data: {
              refreshToken: refreshToken,
            },
          }).then(
            (resp) => {
              localStorage.setItem('accessToken', resp.authenticationResult.accessToken);
              refreshingToken = null;
              return resp;
            },
            (err) => {
              refreshingToken = null;
              toast.error('Please login first');
              history.push('/login');
              return Promise.reject(err);
            },
          );
        return refreshingToken.then(() => {
          const originalRequest = error.config;
          originalRequest._retry = true;
          return Promise.resolve(request(originalRequest.url, originalRequest));
        });
      } else {
        history.push('/login');
      }
    }

    if (status === 428) {
      history.push('/db-migrate');
    }

    if (status === 500) {
      history.push('/offline');
    }

    return Promise.reject(error);
  },
);

export type ReuqestConfig = {
  method?: AxiosRequestConfig['method'];
  data?: unknown;
  timeout?: number;
  signal?: AbortSignal;
  headers?: Record<string, string>;
};

export const request = (path: string, config?: ReuqestConfig) => {
  const { method = 'get', data, timeout, headers, signal } = config || {};

  const cancelTokenSource = axios.CancelToken.source();
  const token = localStorage.getItem('accessToken');
  const h = { ...headers };
  if (token) {
    h.Authorization = `Bearer ${token}`;
  }

  const params: any = {
    url: path,
    method,
    timeout,
    headers: h,
    cancelToken: cancelTokenSource?.token,
  };

  if (['GET', 'get'].includes(method)) {
    params.params = data;
  } else {
    params.data = data;
  }

  const promise = instance.request(params).then((resp) => resp.data);

  if (signal) {
    signal.addEventListener('abort', () => {
      cancelTokenSource?.cancel();
    });
  }

  return promise;
};
