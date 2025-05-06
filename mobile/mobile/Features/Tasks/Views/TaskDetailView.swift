import SwiftUI

struct TaskDetailView: View {
    let task: Task
    @ObservedObject var viewModel: TasksViewModel
    @Environment(\.dismiss) private var dismiss
    
    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 20) {
                HStack {
                    Text(task.title)
                        .font(.title2)
                        .bold()
                    Spacer()
                    PriorityBadge(priority: task.priority)
                }
                
                if let description = task.description {
                    InfoSection(title: "Описание") {
                        Text(description)
                            .foregroundColor(.secondary)
                    }
                }
                
                InfoSection(title: "Статус") {
                    HStack(spacing: 12) {
                        Circle()
                            .fill(task.isCompleted ? Color.green : Color.orange)
                            .frame(width: 12, height: 12)
                        Text(task.status.rawValue)
                            .foregroundColor(.primary)
                            .font(.body)
                        Spacer()
                        Text(task.isCompleted ? "Выполнено" : "В процессе")
                            .font(.subheadline)
                            .foregroundColor(task.isCompleted ? .green : .orange)
                    }
                    .padding(.vertical, 4)
                }
                
                if let deadline = task.deadline {
                    InfoSection(title: "Дедлайн") {
                        HStack {
                            Image(systemName: "calendar")
                                .foregroundColor(.blue)
                            Text(formatDate(deadline))
                                .foregroundColor(.secondary)
                        }
                    }
                }
                
                InfoSection(title: "Создано") {
                    HStack {
                        Image(systemName: "clock")
                            .foregroundColor(.gray)
                        Text(formatDate(task.createdAt))
                            .foregroundColor(.secondary)
                    }
                }
                
                InfoSection(title: "Изменено") {
                    HStack {
                        Image(systemName: "pencil")
                            .foregroundColor(.gray)
                        Text(formatDate(task.updatedAt ?? ""))
                            .foregroundColor(.secondary)
                    }
                }
                
                VStack(spacing: 12) {
                    Button {
                        viewModel.updateTaskStatus(id: task.id, isCompleted: !task.isCompleted)
                    } label: {
                        HStack {
                            Image(systemName: task.isCompleted ? "checkmark.circle.fill" : "circle")
                            Text(task.isCompleted ? "Отметить как невыполненную" : "Отметить как выполненную")
                        }
                        .frame(maxWidth: .infinity)
                    }
                    .buttonStyle(.bordered)
                    .tint(task.isCompleted ? .red : .green)
                    
                    Button(role: .destructive) {
                        viewModel.deleteTask(id: task.id)
                        dismiss()
                    } label: {
                        HStack {
                            Image(systemName: "trash")
                            Text("Удалить задачу")
                        }
                        .frame(maxWidth: .infinity)
                    }
                    .buttonStyle(.bordered)
                }
                .padding(.top, 8)
            }
            .padding()
        }
        .navigationBarTitleDisplayMode(.inline)
    }
    
    // MARK: - Info Section
    private struct InfoSection<Content: View>: View {
        let title: String
        let content: Content
        
        init(title: String, @ViewBuilder content: () -> Content) {
            self.title = title
            self.content = content()
        }
        
        var body: some View {
            VStack(alignment: .leading, spacing: 8) {
                Text(title)
                    .font(.headline)
                    .foregroundColor(.primary)
                content
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
                .background(Self.priorityColor(for: priority))
                .foregroundColor(.white)
                .cornerRadius(4)
        }
        
        private static func priorityColor(for priority: TaskPriority) -> Color {
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
    }
    
    // MARK: - Helper Functions
    private func formatDate(_ dateString: String) -> String {
        let inputFormatter = DateFormatter()
        inputFormatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ss.SSSZ"
        inputFormatter.locale = Locale(identifier: "ru_RU")
        inputFormatter.timeZone = TimeZone(abbreviation: "UTC")
        
        let outputFormatter = DateFormatter()
        outputFormatter.dateFormat = "d MMMM HH:mm"
        outputFormatter.locale = Locale(identifier: "ru_RU")
        outputFormatter.timeZone = TimeZone.current
        
        if let date = inputFormatter.date(from: dateString) {
            return outputFormatter.string(from: date)
        }
        
        inputFormatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ssZ"
        if let date = inputFormatter.date(from: dateString) {
            return outputFormatter.string(from: date)
        }
        
        return dateString
    }
}
