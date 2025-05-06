import SwiftUI

struct CreateTaskView: View {
    @Environment(\.dismiss) private var dismiss
    @ObservedObject var viewModel: TasksViewModel
    
    @State private var title: String = ""
    @State private var description: String = ""
    @State private var deadline: Date = Date()
    @State private var priority: TaskPriority = .medium
    @State private var showDatePicker: Bool = false
    @State private var showDateError: Bool = false
    
    var body: some View {
        NavigationView {
            Form {
                Section("Основная информация") {
                    TextField("Название", text: $title)
                    TextField("Описание", text: $description)
                }
                
                Section("Приоритет") {
                    Picker("Приоритет", selection: $priority) {
                        Text("Низкий").tag(TaskPriority.low)
                        Text("Средний").tag(TaskPriority.medium)
                        Text("Высокий").tag(TaskPriority.high)
                        Text("Критический").tag(TaskPriority.critical)
                    }
                    .pickerStyle(.menu)
                }
                
                Section("Дедлайн") {
                    Toggle("Установить дедлайн", isOn: $showDatePicker)
                    if showDatePicker {
                        DatePicker(
                            "Дедлайн",
                            selection: $deadline,
                            in: Date()...,
                            displayedComponents: [.date, .hourAndMinute]
                        )
                    }
                }
            }
            .navigationTitle("Новая задача")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Отмена") {
                        dismiss()
                    }
                }
                
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Создать") {
                        createTask()
                    }
                    .disabled(title.isEmpty)
                }
            }
            .alert("Ошибка", isPresented: $showDateError) {
                Button("OK", role: .cancel) { }
            } message: {
                Text("Дедлайн не может быть в прошлом")
            }
        }
    }
    
    private func createTask() {
        let deadlineString = showDatePicker ? formatDate(deadline) : nil
        
        if showDatePicker && deadline < Date() {
            showDateError = true
            return
        }
        
        viewModel.createTask(
            title: title,
            description: description.isEmpty ? nil : description,
            deadline: deadlineString,
            priority: priority
        )
        
        dismiss()
    }
    
    private func formatDate(_ date: Date) -> String {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime]
        return formatter.string(from: date)
    }
}
