//
//  NetworkService.swift
//  mobile
//
//  Created by Богдан Тарченко on 05.05.2025.
//

import Foundation
import Combine

enum NetworkError: Error {
    case invalidURL
    case invalidResponse
    case decodingError
    case serverError(String)
}

class NetworkService {
    private let baseURL: String
    private let session: URLSession
    
    init(baseURL: String = "http://localhost:8080", session: URLSession = .shared) {
        self.baseURL = baseURL
        self.session = session
    }
    
    func request<T: Decodable>(_ endpoint: Endpoint) -> AnyPublisher<T, Error> {
        guard let url = URL(string: baseURL + endpoint.path) else {
            return Fail(error: NetworkError.invalidURL).eraseToAnyPublisher()
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = endpoint.method.rawValue
        request.allHTTPHeaderFields = endpoint.headers
        
        if let body = endpoint.body {
            request.httpBody = try? JSONEncoder().encode(body)
        }
        
        return session.dataTaskPublisher(for: request)
            .tryMap { data, response in
                guard let httpResponse = response as? HTTPURLResponse else {
                    throw NetworkError.invalidResponse
                }
                
                guard (200...299).contains(httpResponse.statusCode) else {
                    throw NetworkError.serverError("Status code: \(httpResponse.statusCode)")
                }
                
                return data
            }
            .decode(type: T.self, decoder: JSONDecoder())
            .mapError { error in
                if let decodingError = error as? DecodingError {
                    return NetworkError.decodingError
                }
                return error
            }
            .eraseToAnyPublisher()
    }
}
