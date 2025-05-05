//
//  Endpoint.swift
//  mobile
//
//  Created by Богдан Тарченко on 05.05.2025.
//

enum HTTPMethod: String {
    case get = "GET"
    case post = "POST"
    case put = "PUT"
    case delete = "DELETE"
}

struct Endpoint {
    let path: String
    let method: HTTPMethod
    let headers: [String: String]
    let body: Encodable?
    
    init(
        path: String,
        method: HTTPMethod = .get,
        headers: [String: String] = [:],
        body: Encodable? = nil
    ) {
        self.path = path
        self.method = method
        self.headers = headers
        self.body = body
    }
}

extension Endpoint {
    static func tasks() -> Endpoint {
        Endpoint(path: "/api/tasks")
    }
    
    static func task(id: Int) -> Endpoint {
        Endpoint(path: "/api/tasks/\(id)")
    }
    
    static func createTask(_ task: CreateTaskDTO) -> Endpoint {
        Endpoint(
            path: "/api/tasks",
            method: .post,
            headers: ["Content-Type": "application/json"],
            body: task
        )
    }
    
    static func updateTask(id: Int, _ task: UpdateTaskDTO) -> Endpoint {
        Endpoint(
            path: "/api/tasks/\(id)",
            method: .put,
            headers: ["Content-Type": "application/json"],
            body: task
        )
    }
    
    static func deleteTask(id: Int) -> Endpoint {
        Endpoint(
            path: "/api/tasks/\(id)",
            method: .delete
        )
    }
}
