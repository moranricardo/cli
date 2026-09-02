echo "== go.mod =="; cat go.mod
echo ""; echo "== cmd =="; ls -R cmd
echo ""; echo "== internal =="; ls -R internal | head -n 80
echo ""; echo "== workflows =="; ls -la .github/workflows/; cat .github/workflows/*.yml
echo ""; echo "== README =="; cat README.md
echo ""; echo "== go test =="; go test ./... 2>&1 | tail -n 30
