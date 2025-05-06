import Foundation
import Combine
import SwiftUI

@MainActor
final class TasksViewModel: ObservableObject {
    @Published var tasks: [Task] = []
    @Published var isLoading: Bool = false
    @Published var error: String?
    @Published var currentPage: Int = 1
    @Published var totalPages: Int = 1
    
    private let taskService: TaskService
    private var cancellables = Set<AnyCancellable>()
    
    init(taskService: TaskService = TaskService()) {
        self.taskService = taskService
        fetchTasks()
    }
    
    func fetchTasks(page: Int = 1, pageSize: Int = 999) {
        isLoading = true
        error = nil
        
        taskService.fetchTasks(page: page, pageSize: pageSize)
            .receive(on: DispatchQueue.main)
            .sink { [weak self] completion in
                self?.isLoading = false
                if case .failure(let error) = completion {
                    self?.error = error.localizedDescription
                }
            } receiveValue: { [weak self] response in
                self?.tasks = response.items
                self?.totalPages = response.meta.totalPages
            }
            .store(in: &cancellables)
    }
    
    func loadNextPage() {
        guard currentPage < totalPages, !isLoading else { return }
        currentPage += 1
        fetchTasks()
    }
    
    func refreshTasks() {
        currentPage = 1
        fetchTasks()
    }
    
    func updateTaskStatus(id: String, isCompleted: Bool) {
        let dto = UpdateTaskStatusDTO(isCompleted: isCompleted)
        
        taskService.updateTaskStatus(id: id, dto: dto)
            .receive(on: DispatchQueue.main)
            .sink { [weak self] completion in
                if case .failure(let error) = completion {
                    self?.error = error.localizedDescription
                }
            } receiveValue: { [weak self] updatedTask in
                if let index = self?.tasks.firstIndex(where: { $0.id == updatedTask.id }) {
                    self?.tasks[index] = updatedTask
                }
            }
            .store(in: &cancellables)
    }
    
    func deleteTask(id: String) {
        withAnimation {
            tasks.removeAll { $0.id == id }
        }
        
        taskService.deleteTask(id: id)
            .receive(on: DispatchQueue.main)
            .sink { [weak self] completion in
                if case .failure(let error) = completion {
                    self?.error = error.localizedDescription
                    self?.fetchTasks()
                }
            } receiveValue: { _ in }
            .store(in: &cancellables)
    }
    
    func createTask(title: String, description: String?, deadline: String?, priority: TaskPriority) {
        let dto = CreateTaskDTO(
            title: title,
            description: description,
            deadline: deadline,
            priority: priority
        )
        
        taskService.createTask(dto)
            .receive(on: DispatchQueue.main)
            .sink { [weak self] completion in
                if case .failure(let error) = completion {
                    self?.error = error.localizedDescription
                }
            } receiveValue: { [weak self] newTask in
                self?.tasks.insert(newTask, at: 0)
            }
            .store(in: &cancellables)
    }
}
