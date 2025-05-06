import SwiftUI

struct TasksView: View {
    @StateObject private var viewModel = TasksViewModel()
    
    var body: some View {
        NavigationView {
            ZStack {
                if viewModel.isLoading && viewModel.tasks.isEmpty {
                    LoadingView()
                } else if let error = viewModel.error {
                    ErrorView(error: error) {
                        viewModel.refreshTasks()
                    }
                } else if viewModel.tasks.isEmpty {
                    EmptyTasksView()
                } else {
                    TasksListView(viewModel: viewModel)
                }
            }
            .navigationTitle("Задачи")
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
    var body: some View {
        VStack(spacing: 16) {
            Image(systemName: "checklist")
                .font(.system(size: 50))
                .foregroundColor(.gray)
            Text("Нет задач")
                .font(.title2)
                .foregroundColor(.gray)
        }
    }
}

// MARK: - Tasks List View
private struct TasksListView: View {
    @ObservedObject var viewModel: TasksViewModel
    
    var body: some View {
        List {
            ForEach(viewModel.tasks) { task in
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
    outputFormatter.dateFormat = "d MMMM yyyy"
    outputFormatter.locale = Locale(identifier: "ru_RU")
    
    if let date = inputFormatter.date(from: dateString) {
        return outputFormatter.string(from: date)
    }
    return dateString
}

#Preview {
    TasksView()
}
