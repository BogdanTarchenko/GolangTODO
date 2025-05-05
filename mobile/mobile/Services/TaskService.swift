//
//  TaskService.swift
//  mobile
//
//  Created by Богдан Тарченко on 05.05.2025.
//

import Foundation
import Combine

protocol TaskServiceProtocol {
    func getTasks(filter: TaskFilter) -> AnyPublisher<[Task], Error>
    func getTask(id: Int) -> AnyPublisher<Task, Error>
    func createTask(_ task: CreateTaskDTO) -> AnyPublisher<Task, Error>
    func updateTask(id: Int, _ task: UpdateTaskDTO) -> AnyPublisher<Task, Error>
    func deleteTask(id: Int) -> AnyPublisher<Void, Error>
    func updateTaskStatus(id: Int, status: TaskStatus) -> AnyPublisher<Task, Error>
}

class TaskService: TaskServiceProtocol {
    private let networkService: NetworkService
    
    init(networkService: NetworkService = NetworkService()) {
        self.networkService = networkService
    }
    
    func getTasks(filter: TaskFilter) -> AnyPublisher<[Task], Error> {
        var path = "/api/tasks"
        var queryItems: [URLQueryItem] = []
        
        if let status = filter.status {
            queryItems.append(URLQueryItem(name: "status", value: status.rawValue))
        }
        if let priority = filter.priority {
            queryItems.append(URLQueryItem(name: "priority", value: priority.rawValue))
        }
        queryItems.append(URLQueryItem(name: "page", value: String(filter.page)))
        queryItems.append(URLQueryItem(name: "page_size", value: String(filter.pageSize)))
        
        if !queryItems.isEmpty {
            path += "?" + queryItems.map { "\($0.name)=\($0.value ?? "")" }.joined(separator: "&")
        }
        
        let endpoint = Endpoint(path: path)
        return networkService.request(endpoint)
    }
    
    func getTask(id: Int) -> AnyPublisher<Task, Error> {
        let endpoint = Endpoint.task(id: id)
        return networkService.request(endpoint)
    }
    
    func createTask(_ task: CreateTaskDTO) -> AnyPublisher<Task, Error> {
        let endpoint = Endpoint.createTask(task)
        return networkService.request(endpoint)
    }
    
    func updateTask(id: Int, _ task: UpdateTaskDTO) -> AnyPublisher<Task, Error> {
        let endpoint = Endpoint.updateTask(id: id, task)
        return networkService.request(endpoint)
    }
    
    func deleteTask(id: Int) -> AnyPublisher<Void, Error> {
        let endpoint = Endpoint.deleteTask(id: id)
        return networkService.request(endpoint)
            .map { (_: EmptyResponse) in () }
            .eraseToAnyPublisher()
    }
    
    func updateTaskStatus(id: Int, status: TaskStatus) -> AnyPublisher<Task, Error> {
        let updateDTO = UpdateTaskDTO(status: status)
        return updateTask(id: id, updateDTO)
    }
}

// MARK: - Error Handling
extension TaskService {
    enum TaskError: LocalizedError {
        case taskNotFound
        case invalidData
        case networkError
        
        var errorDescription: String? {
            switch self {
            case .taskNotFound:
                return "Task not found"
            case .invalidData:
                return "Invalid data"
            case .networkError:
                return "Network error"
            }
        }
    }
}

// MARK: - Date Coding
extension TaskService {
    static let dateFormatter: DateFormatter = {
        let formatter = DateFormatter()
        formatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ssZ"
        return formatter
    }()
    
    static let jsonDecoder: JSONDecoder = {
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .formatted(dateFormatter)
        return decoder
    }()
    
    static let jsonEncoder: JSONEncoder = {
        let encoder = JSONEncoder()
        encoder.dateEncodingStrategy = .formatted(dateFormatter)
        return encoder
    }()
}

// MARK: - Empty Response
private struct EmptyResponse: Decodable {}
