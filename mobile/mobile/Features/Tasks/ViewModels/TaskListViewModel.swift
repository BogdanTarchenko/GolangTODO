//
//  TaskListViewModel.swift
//  mobile
//
//  Created by Богдан Тарченко on 05.05.2025.
//

import SwiftUI
import Foundation
import Combine

class TaskListViewModel: ObservableObject {
    // MARK: - Published Properties
    @Published var tasks: [Task] = []
    @Published var isLoading = false
    @Published var error: String?
    @Published var filter = TaskFilter()
    
    // MARK: - Private Properties
    private let taskService: TaskServiceProtocol
    private var cancellables = Set<AnyCancellable>()
    
    // MARK: - Initialization
    init(taskService: TaskServiceProtocol = TaskService()) {
        self.taskService = taskService
    }
    
    // MARK: - Public Methods
    func loadTasks() {
        isLoading = true
        error = nil
        
        taskService.getTasks(filter: filter)
            .receive(on: DispatchQueue.main)
            .sink { [weak self] completion in
                self?.isLoading = false
                if case .failure(let error) = completion {
                    self?.error = error.localizedDescription
                }
            } receiveValue: { [weak self] tasks in
                self?.tasks = tasks
            }
            .store(in: &cancellables)
    }
    
    func createTask(title: String, description: String?, deadline: Date?, priority: TaskPriority?) {
        let taskDTO = CreateTaskDTO(
            title: title,
            description: description,
            deadline: deadline,
            priority: priority
        )
        
        taskService.createTask(taskDTO)
            .receive(on: DispatchQueue.main)
            .sink { [weak self] completion in
                if case .failure(let error) = completion {
                    self?.error = error.localizedDescription
                }
            } receiveValue: { [weak self] _ in
                self?.loadTasks()
            }
            .store(in: &cancellables)
    }
    
    func updateTaskStatus(id: Int, status: TaskStatus) {
        taskService.updateTaskStatus(id: id, status: status)
            .receive(on: DispatchQueue.main)
            .sink { [weak self] completion in
                if case .failure(let error) = completion {
                    self?.error = error.localizedDescription
                }
            } receiveValue: { [weak self] _ in
                self?.loadTasks()
            }
            .store(in: &cancellables)
    }
    
    func deleteTask(id: Int) {
        taskService.deleteTask(id: id)
            .receive(on: DispatchQueue.main)
            .sink { [weak self] completion in
                if case .failure(let error) = completion {
                    self?.error = error.localizedDescription
                }
            } receiveValue: { [weak self] _ in
                self?.loadTasks()
            }
            .store(in: &cancellables)
    }
    
    // MARK: - Filter Methods
    func updateFilter(status: TaskStatus?) {
        filter.status = status
        loadTasks()
    }
    
    func updateFilter(priority: TaskPriority?) {
        filter.priority = priority
        loadTasks()
    }
    
    func updatePage(_ page: Int) {
        filter.page = page
        loadTasks()
    }
    
    // MARK: - Helper Methods
    func taskColor(for task: Task) -> Color {
        if task.status == .completed {
            return .gray
        }
        
        if let deadline = task.deadline {
            let daysUntilDeadline = Calendar.current.dateComponents([.day], from: Date(), to: deadline).day ?? 0
            
            if daysUntilDeadline < 0 {
                return .red
            } else if daysUntilDeadline < 3 {
                return .orange
            }
        }
        
        return .primary
    }
}
