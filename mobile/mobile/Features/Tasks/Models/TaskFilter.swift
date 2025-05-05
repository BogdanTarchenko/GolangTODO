//
//  TaskFilter.swift
//  mobile
//
//  Created by Богдан Тарченко on 05.05.2025.
//

struct TaskFilter {
    var status: TaskStatus?
    var priority: TaskPriority?
    var page: Int = 1
    var pageSize: Int = 10
}
