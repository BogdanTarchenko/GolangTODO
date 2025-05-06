import SwiftUI

struct TaskFiltersView: View {
    @Environment(\.dismiss) private var dismiss
    @ObservedObject var viewModel: TasksViewModel
    
    var body: some View {
        NavigationView {
            Form {
                statusSection
                prioritySection
                sortSection
            }
            .navigationTitle("Фильтры")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Сбросить") {
                        viewModel.resetFilters()
                    }
                }
                
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Применить") {
                        viewModel.fetchTasks()
                        dismiss()
                    }
                }
            }
        }
    }
    
    private var statusSection: some View {
        Section("Статус") {
            Picker("Статус", selection: $viewModel.selectedStatus) {
                Text("Все").tag(nil as TaskStatus?)
                Text("Активно").tag(TaskStatus.active)
                Text("Выполнено").tag(TaskStatus.completed)
                Text("Просрочено").tag(TaskStatus.overdue)
                Text("Сделано с опозданием").tag(TaskStatus.late)
            }
        }
    }
    
    private var prioritySection: some View {
        Section("Приоритет") {
            Picker("Приоритет", selection: $viewModel.selectedPriority) {
                Text("Все").tag(nil as TaskPriority?)
                Text("Низкий").tag(TaskPriority.low)
                Text("Средний").tag(TaskPriority.medium)
                Text("Высокий").tag(TaskPriority.high)
                Text("Критический").tag(TaskPriority.critical)
            }
        }
    }
    
    private var sortSection: some View {
        Section("Сортировка") {
            Picker("Сортировать по", selection: $viewModel.sortBy) {
                Text("Все").tag(nil as TaskSortField?)
                Text("Дата создания").tag(TaskSortField.createdAt)
                Text("Дедлайн").tag(TaskSortField.deadline)
                Text("Приоритет").tag(TaskSortField.priority)
            }
            
            if viewModel.sortBy != nil {
                Picker("Порядок", selection: $viewModel.sortOrder) {
                    Text("По возрастанию").tag(SortOrder.asc)
                    Text("По убыванию").tag(SortOrder.desc)
                }
            }
        }
    }
}
