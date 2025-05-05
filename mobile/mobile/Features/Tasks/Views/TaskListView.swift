//
//  TaskListView.swift
//  mobile
//
//  Created by Богдан Тарченко on 05.05.2025.
//

import SwiftUI

struct TaskListView: View {
    @StateObject private var viewModel = TaskListViewModel()
    @State private var showingAddTask = false
    @State private var showingFilters = false
    
    var body: some View {
        NavigationView {
            ZStack {
                if viewModel.isLoading {
                    ProgressView()
                } else {
                    taskList
                }
            }
            .navigationTitle("Задачи")
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button(action: { showingAddTask = true }) {
                        Image(systemName: "plus")
                    }
                }
                
                ToolbarItem(placement: .navigationBarLeading) {
                    Button(action: { showingFilters = true }) {
                        Image(systemName: "line.3.horizontal.decrease.circle")
                    }
                }
            }
            .sheet(isPresented: $showingAddTask) {
                TaskFormView(viewModel: viewModel)
            }
            .sheet(isPresented: $showingFilters) {
                TaskFilterView(viewModel: viewModel)
            }
            .alert("Ошибка", isPresented: .constant(viewModel.error != nil)) {
                Button("OK") {
                    viewModel.error = nil
                }
            } message: {
                Text(viewModel.error ?? "")
            }
        }
        .onAppear {
            viewModel.loadTasks()
        }
    }
    
    private var taskList: some View {
        List {
            ForEach(viewModel.tasks) { task in
                TaskRowView(task: task, viewModel: viewModel)
            }
            .onDelete { indexSet in
                for index in indexSet {
                    viewModel.deleteTask(id: viewModel.tasks[index].id)
                }
            }
        }
        .refreshable {
            viewModel.loadTasks()
        }
    }
}

struct TaskRowView: View {
    let task: Task
    let viewModel: TaskListViewModel
    
    var body: some View {
        HStack {
            VStack(alignment: .leading, spacing: 4) {
                Text(task.title)
                    .font(.headline)
                    .foregroundColor(viewModel.taskColor(for: task))
                
                if let description = task.description {
                    Text(description)
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                }
                
                HStack {
                    if let deadline = task.deadline {
                        Label(deadline.formatted(date: .abbreviated, time: .omitted),
                              systemImage: "calendar")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    
                    Spacer()
                    
                    Text(task.priority.rawValue)
                        .font(.caption)
                        .padding(.horizontal, 8)
                        .padding(.vertical, 4)
                        .background(priorityColor(for: task.priority))
                        .foregroundColor(.white)
                        .cornerRadius(4)
                }
            }
            
            Spacer()
            
            Button(action: {
                let newStatus: TaskStatus = task.status == .completed ? .active : .completed
                viewModel.updateTaskStatus(id: task.id, status: newStatus)
            }) {
                Image(systemName: task.status == .completed ? "checkmark.circle.fill" : "circle")
                    .foregroundColor(task.status == .completed ? .green : .gray)
            }
        }
        .padding(.vertical, 4)
    }
    
    private func priorityColor(for priority: TaskPriority) -> Color {
        switch priority {
        case .critical:
            return .red
        case .high:
            return .orange
        case .medium:
            return .blue
        case .low:
            return .green
        }
    }
}

struct TaskFilterView: View {
    @Environment(\.dismiss) private var dismiss
    let viewModel: TaskListViewModel
    
    var body: some View {
        NavigationView {
            Form {
                Section("Статус") {
                    Picker("Статус", selection: Binding(
                        get: { viewModel.filter.status },
                        set: { viewModel.updateFilter(status: $0) }
                    )) {
                        Text("Все").tag(Optional<TaskStatus>.none)
                        Text("Активные").tag(Optional<TaskStatus>.some(.active))
                        Text("Завершенные").tag(Optional<TaskStatus>.some(.completed))
                        Text("Просроченные").tag(Optional<TaskStatus>.some(.overdue))
                    }
                }
                
                Section("Приоритет") {
                    Picker("Приоритет", selection: Binding(
                        get: { viewModel.filter.priority },
                        set: { viewModel.updateFilter(priority: $0) }
                    )) {
                        Text("Все").tag(Optional<TaskPriority>.none)
                        Text("Критический").tag(Optional<TaskPriority>.some(.critical))
                        Text("Высокий").tag(Optional<TaskPriority>.some(.high))
                        Text("Средний").tag(Optional<TaskPriority>.some(.medium))
                        Text("Низкий").tag(Optional<TaskPriority>.some(.low))
                    }
                }
            }
            .navigationTitle("Фильтры")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Готово") {
                        dismiss()
                    }
                }
            }
        }
    }
}

#Preview {
    TaskListView()
}
