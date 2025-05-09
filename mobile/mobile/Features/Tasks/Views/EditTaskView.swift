import SwiftUI

struct EditTaskView: View {
    @Environment(\.dismiss) private var dismiss
    @ObservedObject var viewModel: TasksViewModel
    
    let task: Task
    @State private var title: String
    @State private var description: String
    @State private var priority: TaskPriority
    @State private var deadline: Date?
    @State private var deadlineError: String?
    @State private var titleError: String?
    
    init(task: Task, viewModel: TasksViewModel) {
        self.task = task
        self.viewModel = viewModel
        _title = State(initialValue: task.title)
        _description = State(initialValue: task.description ?? "")
        _priority = State(initialValue: task.priority)
        if let deadlineString = task.deadline {
            let formatter = DateFormatter()
            formatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ssZ"
            _deadline = State(initialValue: formatter.date(from: deadlineString))
        } else {
            _deadline = State(initialValue: Calendar.current.date(byAdding: .day, value: 1, to: Date()))
        }
    }
    
    var body: some View {
        NavigationView {
            Form {
                mainInfoSection
                prioritySection
                deadlineSection
            }
            .navigationTitle("Редактирование задачи")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Отмена") {
                        dismiss()
                    }
                }
                
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Сохранить") {
                        saveTask()
                    }
                    .disabled(title.isEmpty)
                }
            }
        }
    }
    
    private var mainInfoSection: some View {
        Section("Основная информация") {
            TextField("Название", text: $title)
                .onChange(of: title) { _ in
                    titleError = nil
                }
            if let error = titleError {
                Text(error)
                    .foregroundColor(.red)
                    .font(.caption)
            }
            TextEditor(text: $description)
                .frame(height: 100)
        }
    }
    
    private var prioritySection: some View {
        Section("Приоритет") {
            Picker("Приоритет", selection: $priority) {
                Text("Низкий").tag(TaskPriority.low)
                Text("Средний").tag(TaskPriority.medium)
                Text("Высокий").tag(TaskPriority.high)
                Text("Критический").tag(TaskPriority.critical)
            }
        }
    }
    
    private var deadlineSection: some View {
        Section("Дедлайн") {
            Toggle("Установить дедлайн", isOn: Binding(
                get: { deadline != nil },
                set: { isOn in
                    if isOn {
                        deadline = Calendar.current.date(byAdding: .day, value: 1, to: Date())
                    } else {
                        deadline = nil
                    }
                    deadlineError = nil
                }
            ))
            
            if deadline != nil {
                HStack {
                    DatePicker(
                        "Дедлайн",
                        selection: Binding(
                            get: { deadline ?? Date() },
                            set: { newDate in
                                if newDate > Date() {
                                    deadline = newDate
                                    deadlineError = nil
                                } else {
                                    deadlineError = "Дедлайн не может быть в прошлом"
                                }
                            }
                        ),
                        in: Date()...,
                        displayedComponents: [.date, .hourAndMinute]
                    )
                    
                    Button {
                        deadline = nil
                        deadlineError = nil
                    } label: {
                        Image(systemName: "xmark.circle.fill")
                            .foregroundColor(.gray)
                    }
                }
                
                if let error = deadlineError {
                    Text(error)
                        .foregroundColor(.red)
                        .font(.caption)
                }
            }
        }
    }
    
    private func saveTask() {
        if title.count < 4 {
            titleError = "Название должно содержать минимум 4 символа"
            return
        }
        
        let formatter = DateFormatter()
        formatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ss'Z'"
        formatter.timeZone = TimeZone(abbreviation: "UTC")
        
        if let deadline = deadline, deadline <= Date() {
            deadlineError = "Дедлайн не может быть в прошлом"
            return
        }
        
        let dto = UpdateTaskDTO(
            title: title,
            description: description.isEmpty ? nil : description,
            deadline: deadline.map { formatter.string(from: $0) },
            priority: priority
        )
        
        viewModel.updateTask(id: task.id, dto: dto)
        dismiss()
    }
}
