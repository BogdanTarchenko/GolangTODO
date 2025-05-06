import Foundation
import Combine
import SwiftUI

@MainActor
final class TasksViewModel: ObservableObject {
    @Published var tasks: [Task]?
    @Published var isLoading: Bool = false
    @Published var error: String?
    @Published var currentPage: Int = 1
    @Published var totalPages: Int = 1
    
    @Published var selectedStatus: TaskStatus? {
        didSet { 
            currentPage = 1
            fetchTasks() 
        }
    }
    @Published var selectedPriority: TaskPriority? {
        didSet { 
            currentPage = 1
            fetchTasks() 
        }
    }
    @Published var sortBy: TaskSortField? {
        didSet { 
            currentPage = 1
            fetchTasks() 
        }
    }
    @Published var sortOrder: SortOrder? {
        didSet { 
            currentPage = 1
            fetchTasks() 
        }
    }
    
    private let taskService: TaskService
    private var cancellables = Set<AnyCancellable>()
    
    init(taskService: TaskService = TaskService()) {
        self.taskService = taskService
        fetchTasks()
    }
    
    func fetchTasks() {
        isLoading = true
        error = nil
        
        taskService.fetchTasks(
            status: selectedStatus,
            priority: selectedPriority,
            sortBy: sortBy?.rawValue,
            sortOrder: sortOrder?.rawValue
        )
        .receive(on: DispatchQueue.main)
        .sink { [weak self] completion in
            self?.isLoading = false
            if case .failure(let error) = completion {
                self?.error = error.localizedDescription
            }
        } receiveValue: { [weak self] response in
            if self?.currentPage == 1 {
                self?.tasks = response.items ?? []
            } else {
                if let currentTasks = self?.tasks, let newItems = response.items {
                    self?.tasks = currentTasks + newItems
                }
            }
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
                if var currentTasks = self?.tasks,
                   let index = currentTasks.firstIndex(where: { $0.id == updatedTask.id }) {
                    currentTasks[index] = updatedTask
                    self?.tasks = currentTasks
                }
            }
            .store(in: &cancellables)
    }
    
    func deleteTask(id: String) {
        withAnimation {
            if var currentTasks = tasks {
                currentTasks.removeAll { $0.id == id }
                tasks = currentTasks
            }
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
                if var currentTasks = self?.tasks {
                    currentTasks.insert(newTask, at: 0)
                    self?.tasks = currentTasks
                } else {
                    self?.tasks = [newTask]
                }
            }
            .store(in: &cancellables)
    }
    
    func resetFilters() {
        selectedStatus = nil
        selectedPriority = nil
        sortBy = nil
        sortOrder = nil
        fetchTasks()
    }
    
    var hasActiveFilters: Bool {
        selectedStatus != nil || 
        selectedPriority != nil || 
        sortBy != nil ||
        sortOrder != nil
    }
}
