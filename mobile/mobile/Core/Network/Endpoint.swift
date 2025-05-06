import Foundation
import Alamofire

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
    static func tasks(page: Int = 1, pageSize: Int = 999) -> Endpoint {
        Endpoint(
            path: "/api/tasks?page=\(page)&page_size=\(pageSize)"
        )
    }
    
    static func task(id: String) -> Endpoint {
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
    
    static func updateTask(id: String, _ task: UpdateTaskDTO) -> Endpoint {
        Endpoint(
            path: "/api/tasks/\(id)",
            method: .patch,
            headers: ["Content-Type": "application/json"],
            body: task
        )
    }
    
    static func deleteTask(id: String) -> Endpoint {
        Endpoint(
            path: "/api/tasks/\(id)",
            method: .delete
        )
    }
    
    static func updateTaskStatus(id: String, _ dto: UpdateTaskStatusDTO) -> Endpoint {
        Endpoint(
            path: "/api/tasks/\(id)/status",
            method: .patch,
            headers: ["Content-Type": "application/json"],
            body: dto
        )
    }
}
