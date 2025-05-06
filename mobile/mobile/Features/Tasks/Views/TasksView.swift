import SwiftUI

struct TasksView: View {
    @StateObject private var viewModel = TasksViewModel()
    @State private var showingCreateTask = false
    @State private var showingFilters = false
    
    var body: some View {
        NavigationView {
            ZStack {
                if viewModel.isLoading && viewModel.tasks == nil {
                    LoadingView()
                } else if viewModel.tasks?.isEmpty ?? true {
                    EmptyTasksView(hasFilters: viewModel.hasActiveFilters)
                } else {
                    TasksListView(viewModel: viewModel)
                }
            }
            .navigationTitle("Задачи")
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button {
                        showingFilters = true
                    } label: {
                        Image(systemName: "line.3.horizontal.decrease.circle")
                    }
                }
                
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button {
                        showingCreateTask = true
                    } label: {
                        Image(systemName: "plus")
                    }
                }
            }
            .sheet(isPresented: $showingCreateTask) {
                CreateTaskView(viewModel: viewModel)
            }
            .sheet(isPresented: $showingFilters) {
                TaskFiltersView(viewModel: viewModel)
            }
        }
        .onAppear {
            viewModel.fetchTasks()
        }
    }
}

// MARK: - Loading View
private struct LoadingView: View {
    var body: some View {
        ProgressView("Загрузка задач...")
    }
}

// MARK: - Error View
private struct ErrorView: View {
    let error: String
    let retryAction: () -> Void
    
    var body: some View {
        VStack(spacing: 16) {
            Text("Ошибка: \(error)")
                .foregroundColor(.red)
                .multilineTextAlignment(.center)
            Button("Повторить", action: retryAction)
                .buttonStyle(.bordered)
        }
        .padding()
    }
}

// MARK: - Empty View
private struct EmptyTasksView: View {
    let hasFilters: Bool
    
    var body: some View {
        VStack(spacing: 12) {
            Image(systemName: "tray")
                .font(.system(size: 48))
                .foregroundColor(.secondary)
            
            Text(hasFilters ? "Задач с такими фильтрами не найдено" : "У вас пока нет задач")
                .font(.headline)
                .foregroundColor(.secondary)
            
            if !hasFilters {
                Text("Нажмите + чтобы создать новую задачу")
                    .font(.subheadline)
                    .foregroundColor(.secondary)
            }
        }
    }
}

// MARK: - Tasks List View
private struct TasksListView: View {
    @ObservedObject var viewModel: TasksViewModel
    
    var body: some View {
        List {
            ForEach(viewModel.tasks ?? []) { task in
                TaskRowView(task: task, viewModel: viewModel)
                    .id(task.id)
            }
        }
        .refreshable {
            viewModel.refreshTasks()
        }
    }
}

// MARK: - Task Row View
private struct TaskRowView: View {
    let task: Task
    @ObservedObject var viewModel: TasksViewModel
    
    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            HStack {
                Text(task.title)
                    .font(.headline)
                Spacer()
                PriorityBadge(priority: task.priority)
            }
            
            if let description = task.description {
                Text(description)
                    .font(.subheadline)
                    .foregroundColor(.secondary)
            }
            
            HStack {
                Text("Статус: \(task.status.rawValue)")
                    .font(.caption)
                    .foregroundColor(.blue)
                Spacer()
                if let deadline = task.deadline {
                    Text("Дедлайн: \(formatDate(deadline))")
                        .font(.caption2)
                        .foregroundColor(.red)
                }
            }
        }
        .padding(.vertical, 4)
        .swipeActions(edge: .trailing) {
            Button(role: .destructive) {
                withAnimation {
                    viewModel.deleteTask(id: task.id)
                }
            } label: {
                Label("Удалить", systemImage: "trash")
            }
        }
    }
}

// MARK: - Priority Badge
private struct PriorityBadge: View {
    let priority: TaskPriority
    
    var body: some View {
        Text(priority.rawValue)
            .font(.caption)
            .padding(.horizontal, 8)
            .padding(.vertical, 4)
            .background(priorityColor(for: priority))
            .foregroundColor(.white)
            .cornerRadius(4)
    }
}

// MARK: - Helper Functions
private func priorityColor(for priority: TaskPriority) -> Color {
    switch priority {
    case .low:
        return .blue
    case .medium:
        return .orange
    case .high:
        return .red
    case .critical:
        return .purple
    }
}

private func formatDate(_ dateString: String) -> String {
    let inputFormatter = DateFormatter()
    inputFormatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ssZ"
    inputFormatter.locale = Locale(identifier: "ru_RU")
    
    let outputFormatter = DateFormatter()
    outputFormatter.dateFormat = "d MMMM yyyy HH:mm"
    outputFormatter.locale = Locale(identifier: "ru_RU")
    
    if let date = inputFormatter.date(from: dateString) {
        return outputFormatter.string(from: date)
    }
    return dateString
}

#Preview {
    TasksView()
}
