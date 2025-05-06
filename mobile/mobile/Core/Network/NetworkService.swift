import Foundation
import Combine
import Alamofire

enum NetworkError: Error, CustomStringConvertible {
    case invalidURL
    case invalidResponse
    case decodingError(Error)
    case serverError(Int, Data?)
    case unknown(Error)
    
    var description: String {
        switch self {
        case .invalidURL:
            return "Invalid URL"
        case .invalidResponse:
            return "Invalid response type"
        case .decodingError(let error):
            return "Decoding error: \(error)"
        case .serverError(let statusCode, let data):
            let body = data.flatMap { String(data: $0, encoding: .utf8) } ?? "No response body"
            return "Server error. Status code: \(statusCode). Response body: \(body)"
        case .unknown(let error):
            return "Unknown error: \(error)"
        }
    }
}

class NetworkService {
    private let baseURL: String
    private let decoder: JSONDecoder

    init(baseURL: String = "http://localhost:8080") {
        self.baseURL = baseURL
        self.decoder = JSONDecoder()
        self.decoder.dateDecodingStrategy = .iso8601
    }

    func request<T: Decodable>(_ endpoint: Endpoint) -> AnyPublisher<T, Error> {
        guard let url = URL(string: baseURL + endpoint.path) else {
            return Fail(error: NetworkError.invalidURL).eraseToAnyPublisher()
        }

        let method = HTTPMethod(rawValue: endpoint.method.rawValue)
        let headers = HTTPHeaders(endpoint.headers)
        let encoder = JSONEncoder()
        encoder.dateEncodingStrategy = .iso8601

        let request: DataRequest
        if let body = endpoint.body {
            do {
                let data = try encoder.encode(AnyEncodable(body))
                request = AF.request(url, method: method, parameters: nil, encoding: JSONDataEncoding(data: data), headers: headers)
            } catch {
                return Fail(error: NetworkError.decodingError(error)).eraseToAnyPublisher()
            }
        } else {
            request = AF.request(url, method: method, headers: headers)
        }

        return request
            .validate()
            .publishData()
            .tryMap { response in
                guard let data = response.data else {
                    throw NetworkError.invalidResponse
                }
                if let statusCode = response.response?.statusCode, !(200...299).contains(statusCode) {
                    throw NetworkError.serverError(statusCode, data)
                }
                return data
            }
            .decode(type: T.self, decoder: decoder)
            .mapError { error in
                if let afError = error as? AFError {
                    return NetworkError.unknown(afError)
                } else if let decodingError = error as? DecodingError {
                    return NetworkError.decodingError(decodingError)
                }
                return error
            }
            .eraseToAnyPublisher()
    }
}

struct AnyEncodable: Encodable {
    private let encodeFunc: (Encoder) throws -> Void

    init<T: Encodable>(_ wrapped: T) {
        self.encodeFunc = wrapped.encode
    }

    func encode(to encoder: Encoder) throws {
        try encodeFunc(encoder)
    }
}

struct JSONDataEncoding: ParameterEncoding {
    let data: Data

    func encode(_ urlRequest: URLRequestConvertible, with parameters: Parameters?) throws -> URLRequest {
        var request = try urlRequest.asURLRequest()
        request.httpBody = data
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        return request
    }
}
