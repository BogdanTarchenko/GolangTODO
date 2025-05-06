import Foundation
import Combine

final class TaskService {
    private let networkService: NetworkService
    
    init(networkService: NetworkService = NetworkService()) {
        self.networkService = networkService
    }
    
    func fetchTasks(page: Int = 1, pageSize: Int = 999) -> AnyPublisher<PaginatedTasksDTO, Error> {
        networkService.request(.tasks(page: page, pageSize: pageSize))
    }
    
    func fetchAllTasks() -> AnyPublisher<[Task], Error> {
        networkService.request(.tasks())
            .map { (response: PaginatedTasksDTO) in response.items }
            .eraseToAnyPublisher()
    }
    
    func fetchTask(id: String) -> AnyPublisher<Task, Error> {
        networkService.request(.task(id: id))
    }
    
    func createTask(_ dto: CreateTaskDTO) -> AnyPublisher<Task, Error> {
        networkService.request(.createTask(dto))
    }
    
    func updateTask(id: String, dto: UpdateTaskDTO) -> AnyPublisher<Task, Error> {
        networkService.request(.updateTask(id: id, dto))
    }
    
    func updateTaskStatus(id: String, dto: UpdateTaskStatusDTO) -> AnyPublisher<Task, Error> {
        networkService.request(.updateTaskStatus(id: id, dto))
    }
    
    func deleteTask(id: String) -> AnyPublisher<Void, Error> {
        networkService.request(.deleteTask(id: id))
            .map { (_: EmptyResponse) in () }
            .eraseToAnyPublisher()
    }
}
