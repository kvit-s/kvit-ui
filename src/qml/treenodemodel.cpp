// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#include "treenodemodel.h"

#include <QtQml/qqmlinfo.h>

namespace KvitUi {

namespace {

// The role a delegate reads to get at whatever else the entry said. Past the
// last role Qt defines for itself.
constexpr int EntryRole = Qt::UserRole + 1;

}   // namespace

TreeNodeModel::TreeNodeModel(QObject *parent)
    : QAbstractItemModel(parent)
{
}

TreeNodeModel::~TreeNodeModel() = default;

void TreeNodeModel::setNodes(const QVariantList &nodes)
{
    if (m_nodes == nodes)
        return;

    // A whole-model reset rather than per-row signals: the property is set
    // once with the finished hierarchy, and a diff of two nested lists would
    // be more code than anything in this repository asks for.
    beginResetModel();
    m_nodes = nodes;
    m_root.children.clear();
    build(m_nodes, &m_root);
    endResetModel();

    emit nodesChanged();
}

void TreeNodeModel::build(const QVariantList &entries, Node *parent)
{
    int row = 0;
    for (const QVariant &entry : entries) {
        auto node = std::make_unique<Node>();
        node->parent = parent;
        node->row = row++;
        node->entry = entry;

        if (entry.typeId() == QMetaType::QVariantMap) {
            const QVariantMap map = entry.toMap();
            node->label = map.value(QStringLiteral("label")).toString();
            if (node->label.isEmpty()) {
                qmlWarning(this) << "a tree node has no label";
                node->label = tr("(no label)");
            }
            build(map.value(QStringLiteral("children")).toList(), node.get());
        } else {
            node->label = entry.toString();
        }

        parent->children.push_back(std::move(node));
    }
}

const TreeNodeModel::Node *TreeNodeModel::nodeFor(const QModelIndex &index) const
{
    if (!index.isValid())
        return &m_root;
    return static_cast<const Node *>(index.constInternalPointer());
}

QModelIndex TreeNodeModel::index(int row, int column, const QModelIndex &parent) const
{
    if (!hasIndex(row, column, parent))
        return QModelIndex();

    const Node *node = nodeFor(parent);
    return createIndex(row, column, node->children[size_t(row)].get());
}

QModelIndex TreeNodeModel::parent(const QModelIndex &child) const
{
    if (!child.isValid())
        return QModelIndex();

    const Node *node = nodeFor(child);
    if (!node->parent || node->parent == &m_root)
        return QModelIndex();

    return createIndex(node->parent->row, 0, node->parent);
}

int TreeNodeModel::rowCount(const QModelIndex &parent) const
{
    // A tree model has one column, and Qt asks for the row count of an index
    // in every column: answering for column 1 would give a node children it
    // does not have.
    if (parent.isValid() && parent.column() > 0)
        return 0;
    return int(nodeFor(parent)->children.size());
}

int TreeNodeModel::columnCount(const QModelIndex &) const
{
    return 1;
}

QVariant TreeNodeModel::data(const QModelIndex &index, int role) const
{
    if (!index.isValid())
        return QVariant();

    const Node *node = nodeFor(index);
    switch (role) {
    case Qt::DisplayRole:
        return node->label;
    case EntryRole:
        return node->entry;
    default:
        return QVariant();
    }
}

QHash<int, QByteArray> TreeNodeModel::roleNames() const
{
    return {
        { Qt::DisplayRole, QByteArrayLiteral("display") },
        { EntryRole, QByteArrayLiteral("entry") },
    };
}

}   // namespace KvitUi
