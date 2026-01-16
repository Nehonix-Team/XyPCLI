#!/bin/bash

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║         XyPCLI - Script de Déploiement                       ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""

# Vérifier que le binaire existe
if [ ! -f "xypcli" ]; then
    echo "❌ Erreur: xypcli n'existe pas dans ce répertoire"
    echo "   Exécutez d'abord: go build -o xypcli"
    exit 1
fi

echo "📦 Binaire trouvé: $(ls -lh xypcli | awk '{print $5}')"
echo ""

# Copier le binaire
TARGET="/home/idevo/.nvm/versions/node/v22.19.0/bin/xyp"
echo "📋 Copie vers: $TARGET"

cp -f xypcli "$TARGET"

if [ $? -eq 0 ]; then
    echo "✅ Déploiement réussi!"
    echo ""
    echo "🧪 Test de la version:"
    xyp --version
    echo ""
    echo "📖 Afficher l'aide complète:"
    echo "   xyp help"
    echo ""
    echo "🚀 Commande de test recommandée:"
    echo "   cd ~/Documents/projects"
    echo "   xyp init --name xhs-testing --port 5627 --mode n \\"
    echo "     --desc \"Test de la nouvelle version\" \\"
    echo "     --alias \"xhs\" --author \"iDevo\""
else
    echo "❌ Échec du déploiement"
    echo "   Essayez manuellement:"
    echo "   cp xypcli $TARGET"
    exit 1
fi
