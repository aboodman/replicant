import SwiftUI

struct Cell: View {
    
    var title: String = "Coffee"
    
    var notes: String = "Medium roast, unground"
    
    var completed: Bool = false
    
    init(item: [String:String]) {
        if let title = item["id"] {
            self.title = title
        }
        if let notes = item["text"] {
            self.notes = notes
        }
    }
    init(title: String, notes: String) {
        self.title = title
        self.notes = notes
    }

    var body: some View {
        Button(action: {}) {
            HStack {
                Image(systemName: completed ? "checkmark.square.fill" : "square")
                    .foregroundColor(completed ? .green : .secondary)
                VStack(alignment: HorizontalAlignment.leading) {
                    Text(title)
                        .font(.headline)
                        .foregroundColor(.primary)
                    if self.notes != "" {
                        Text(notes)
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                    }
                }
            }
        }
    }
}

#if DEBUG
struct Cell_Previews: PreviewProvider {
    static var previews: some View {
        Cell(title:"foo", notes:"bar")
    }
}
#endif
