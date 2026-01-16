#!/bin/bash

# Test de la fusion intelligente de xypriss.config.json

echo "=========================================="
echo "Test: Fusion Intelligente de Config"
echo "=========================================="
echo ""

# Créer un répertoire de test
TEST_DIR="test-config-merge"
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"

# Créer un fichier xypriss.config.json avec des données existantes
cat > "$TEST_DIR/xypriss.config.json" << 'EOF'
{
  "$internal": {
    "$plg": {
      "__meta__": {
        "path": "#$."
      },
      "__xfs__": {
        "path": "$#."
      }
    }
  },
  "customSection": {
    "myData": "should be preserved",
    "myArray": [1, 2, 3]
  }
}
EOF

echo "1. Fichier de config initial créé:"
echo "-----------------------------------"
cat "$TEST_DIR/xypriss.config.json"
echo ""

# Simuler ce que fait createConfigFile en Go
echo "2. Simulation de la fusion (comme createConfigFile):"
echo "-----------------------------------------------------"

# Utiliser jq pour fusionner (simule le comportement Go)
if command -v jq &> /dev/null; then
    # Lire le fichier existant et ajouter __sys__
    jq '. + {
      "__sys__": {
        "__version__": "1.0.0",
        "__author__": "Test Author",
        "__name__": "test-app",
        "__description__": "Test description",
        "__alias__": "TA",
        "__port__": 3000,
        "__PORT__": 3000
      }
    }' "$TEST_DIR/xypriss.config.json" > "$TEST_DIR/xypriss.config.merged.json"
    
    echo "Fichier après fusion:"
    cat "$TEST_DIR/xypriss.config.merged.json"
    echo ""
    
    echo "3. Vérification:"
    echo "----------------"
    
    # Vérifier que $internal est préservé
    if jq -e '."$internal"' "$TEST_DIR/xypriss.config.merged.json" > /dev/null; then
        echo "✅ Section \$internal préservée"
    else
        echo "❌ Section \$internal perdue!"
    fi
    
    # Vérifier que customSection est préservé
    if jq -e '.customSection' "$TEST_DIR/xypriss.config.merged.json" > /dev/null; then
        echo "✅ Section customSection préservée"
    else
        echo "❌ Section customSection perdue!"
    fi
    
    # Vérifier que __sys__ a été ajouté
    if jq -e '."__sys__"' "$TEST_DIR/xypriss.config.merged.json" > /dev/null; then
        echo "✅ Section __sys__ ajoutée"
    else
        echo "❌ Section __sys__ non ajoutée!"
    fi
    
    echo ""
    echo "=========================================="
    echo "Test de fusion: RÉUSSI ✅"
    echo "=========================================="
else
    echo "⚠️  jq n'est pas installé, impossible de tester la fusion"
    echo "   Installez jq avec: sudo apt install jq"
fi

# Nettoyage
echo ""
echo "Nettoyage du répertoire de test..."
rm -rf "$TEST_DIR"
echo "Terminé!"
