// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#ifndef KVIT_UI_TREENODEMODEL_H
#define KVIT_UI_TREENODEMODEL_H

#include <memory>
#include <vector>

#include <QAbstractItemModel>
#include <QString>
#include <QVariantList>
#include <QVariantMap>
#include <QtQml/qqmlregistration.h>

namespace KvitUi {

// A hierarchy written out in QML, for a KvitTree whose shape is fixed.
//
// A TreeView only accepts a QAbstractItemModel, and QML has no tree model of
// its own the way it has ListModel — so before this, using KvitTree at all
// meant writing a C++ model. Half the trees in the estate do not need one:
// kvit-cash's account groups and category hierarchy, and the outline in a
// settings page, are a handful of fixed rows that belong beside the view that
// shows them. The other half — kvit-notes' file tree over a vault — still
// wants a real model, and sets KvitTree's `model` directly.
//
// A node is either a string, which is a leaf with that label, or an object
// with `label` and an optional `children` array:
//
//     nodes: [
//         "Inbox",
//         { "label": "Accounts", "children": ["Checking", "Savings"] }
//     ]
//
// An object with no `label` is a mistake that would otherwise draw a blank
// row, so it warns and shows the placeholder text rather than nothing.
class TreeNodeModel : public QAbstractItemModel
{
    Q_OBJECT
    QML_ELEMENT

    Q_PROPERTY(QVariantList nodes READ nodes WRITE setNodes NOTIFY nodesChanged)

public:
    explicit TreeNodeModel(QObject *parent = nullptr);
    ~TreeNodeModel() override;

    QVariantList nodes() const { return m_nodes; }
    void setNodes(const QVariantList &nodes);

    QModelIndex index(int row, int column,
                      const QModelIndex &parent = QModelIndex()) const override;
    QModelIndex parent(const QModelIndex &child) const override;
    int rowCount(const QModelIndex &parent = QModelIndex()) const override;
    int columnCount(const QModelIndex &parent = QModelIndex()) const override;
    QVariant data(const QModelIndex &index, int role = Qt::DisplayRole) const override;
    QHash<int, QByteArray> roleNames() const override;

signals:
    void nodesChanged();

private:
    struct Node
    {
        QString label;
        // What the entry said, so a delegate that wants more than the label
        // can read the rest of it through the `entry` role.
        QVariant entry;
        Node *parent = nullptr;
        // Where this node sits among its siblings. Kept rather than searched
        // for, because parent() is called for every visible row on every
        // scroll and a search there is a scan of the parent's children.
        int row = 0;
        std::vector<std::unique_ptr<Node>> children;
    };

    // The node an index points at, and the root for an invalid index — which
    // is what a TreeView passes when it asks for the top level.
    const Node *nodeFor(const QModelIndex &index) const;
    void build(const QVariantList &entries, Node *parent);

    QVariantList m_nodes;
    Node m_root;
};

}   // namespace KvitUi

#endif   // KVIT_UI_TREENODEMODEL_H
